package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type parser struct {
	basketIndex    map[string]string
	commonReplacer *strings.Replacer
	syntaxRegex    map[string]*regexp.Regexp
	wg             sync.WaitGroup
}

func newParser(convention *Convention) *parser {
	return &parser{
		basketIndex:    makeBasket(convention),
		commonReplacer: makeReplacementSet(convention),
		syntaxRegex:    makeSyntaxRegexp(convention),
	}
}

func (parser *parser) parse(source, dest string) error {
	// remove the templates dir
	if err := os.RemoveAll(dest); err != nil {
		return err
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}
	for pattern := range parser.basketIndex {
		founds, err := filepath.Glob(filepath.Join(source, pattern))
		if err != nil {
			log.Printf("Cannot glob pattern %s: %v\n", pattern, err)
		}
		for _, f := range founds {
			var relPath = strings.Replace(f, source, "", 1)

			if err := parser.parseBasket(f, filepath.Join(dest, parser.basketIndex[pattern], relPath)); err != nil {
				log.Printf("Cannot copy %s: %v\n", f, err)
			}
		}
	}
	return nil
}

func (parser *parser) parseFile(source string, destination string, path string, info os.FileInfo) error {
	log.Println("Process:", path)
	var relPath = strings.Replace(path, source, "", 1)

	if info.IsDir() {
		dir := strings.ReplaceAll(filepath.Join(destination, relPath), ".", "dot@")
		return os.MkdirAll(dir, os.ModePerm)
	}

	ext := filepath.Ext(path)
	syntaxRegexp := parser.syntaxRegex[ext]
	var data, err = ioutil.ReadFile(filepath.Join(source, relPath))
	if err != nil {
		return err
	}

	var res = string(data)

	var syntaxM syntaxVault
	if syntaxRegexp != nil {
		for _, x := range syntaxRegexp.FindAllStringSubmatch(res, -1) {
			syntaxM = syntaxM.put(fmt.Sprintf("{{ %s }}", x[1]))
			res = strings.Replace(res, x[0], syntaxM.lastInsertedKey(), 1)
		}
	}

	res = parser.commonReplacer.Replace(res)
	for y, x := range syntaxM {
		res = strings.Replace(res, syntaxM.getKeyByIdx(y), x, 1)
	}

	dir, file := filepath.Split(filepath.Join(destination, relPath))
	file += ".tmpl"
	if file[0] == '.' {
		file = "dot@" + file[1:]
	}

	log.Println("Done:", filepath.Join(dir, file))

	dir = strings.ReplaceAll(dir, ".", "dot@")
	_ = os.MkdirAll(dir, os.ModePerm)
	return ioutil.WriteFile(filepath.Join(dir, file), []byte(res), os.ModePerm)
}

func (parser *parser) parseBasket(source, destination string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parser.wg.Add(1)
		go func() {
			defer parser.wg.Done()
			if err := parser.parseFile(source, destination, path, info); err != nil {
				log.Println("err", err)
			}
		}()
		return nil
	})
}
