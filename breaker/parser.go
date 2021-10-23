package breaker

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type breaker struct {
	layerIndex       map[string][]string
	contentReplacer  *strings.Replacer
	filenameReplacer *strings.Replacer
	syntaxRegex      map[string][]*matchRule
	root             string

	wg sync.WaitGroup
}

func newBreaker(root string, convention *Guide) *breaker {
	return &breaker{
		layerIndex:       makeLayerIndex(convention),
		contentReplacer:  makeContentReplacementSet(convention),
		filenameReplacer: makeFilenameReplacementSet(convention),
		syntaxRegex:      makeSyntaxRegexp(convention),
		root:             root,
	}
}

func (breaker *breaker) make(source, dest string) error {
	// remove the templates dir
	if err := os.RemoveAll(dest); err != nil {
		return err
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}
	for layerName, patterns := range breaker.layerIndex {
		for _, pattern := range patterns {
			breaker.wg.Add(1)
			go func(layerName, pattern string) {
				defer breaker.wg.Done()
				breaker.breakLayer(source, dest, layerName, pattern)
			}(layerName, pattern)
		}
	}
	return nil
}

func (breaker *breaker) breakLayer(source, dest, layerName, pattern string) {
	var sourcePattern = filepath.Join(source, pattern)
	founds, err := filepath.Glob(sourcePattern)
	if err != nil {
		log.Printf("Cannot glob pattern %s: %v\n", pattern, err)
		return
	}
	for _, filePath := range founds {
		if err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(source, path)
			if err != nil {
				log.Println("Cannot load relPath", err)
				return err
			}

			var destPath = filepath.Join(dest, layerName)
			if info.IsDir() {
				dir := strings.ReplaceAll(filepath.Join(destPath, relPath), ".", "dot@")
				return os.MkdirAll(dir, os.ModePerm)
			}

			log.Println("Breaking: ", path, "to", destPath)
			if err := breaker.breakFile(source, destPath, relPath); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Println("Walk err", err)
		}
	}
}

func (breaker *breaker) breakFile(source, dest, relPath string) error {
	var sourceFilePath = filepath.Join(source, relPath)
	data, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}
	var (
		res           = string(data)
		ext           = filepath.Ext(relPath)
		syntaxRegexps = breaker.syntaxRegex[ext]
		syntaxV       = make(syntaxVault, 0)
	)

	for _, syntaxRegexp := range syntaxRegexps {
		for _, matched := range syntaxRegexp.FindAllStringSubmatch(res, -1) {
			var matchedSyntax = matched[1]
			var replacement string
			if syntaxRegexp.needEscape {
				replacement = fmt.Sprintf("%s ", matchedSyntax)
			} else {
				replacement = fmt.Sprintf("{{%s }}", matchedSyntax)
			}

			syntaxV = syntaxV.put(replacement)

			var key = syntaxV.lastInsertedKey()
			if syntaxRegexp.newline {
				key += "\n"
			}
			res = strings.Replace(res, matched[0], key, 1)
		}
	}

	res = breaker.contentReplacer.Replace(res)
	if len(syntaxV) > 0 {
		res = syntaxV.replace(res)
	}

	var (
		dirPath, fileName = breaker.makeDestPath(relPath)
		destDirPath       = filepath.Join(dest, dirPath)
		destFilePath      = filepath.Join(destDirPath, fileName)
	)
	_ = os.MkdirAll(destDirPath, os.ModePerm)
	if err := ioutil.WriteFile(destFilePath, []byte(res), os.ModePerm); err != nil {
		return err
	}
	log.Println("Done:", destFilePath)
	return nil
}
func (breaker *breaker) makeDestPath(relPath string) (string, string) {
	dirPath, fileName := filepath.Split(relPath)
	fileName = breaker.filenameReplacer.Replace(fileName)
	fileName = breaker.makeDestFilename(fileName)

	dirPath = breaker.filenameReplacer.Replace(dirPath)
	dirPath = strings.ReplaceAll(dirPath, ".", "dot@")

	return dirPath, fileName
}
func (breaker *breaker) makeDestFilename(fileName string) string {
	fileName += ".tmpl"
	if fileName[0] == '.' {
		fileName = "dot@" + fileName[1:]
	}
	return fileName
}
