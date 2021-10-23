package parser

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"mvdan.cc/gofumpt/format"
)

type Parser interface {
	WithOption(opt ...Opt) Parser
	Parse(source, destination string) error
}

var _ Parser = &parser{}

// parser responsible for parsing all templates in the source folder to destination folder with provider parameter
// custom the destination file name with altName
type parser struct {
	*config
}

func (jw *parser) WithOption(opts ...Opt) Parser {
	var res = parser{defaultConfig()}
	*res.config = *jw.config
	for _, opt := range opts {
		opt(res.config)
	}
	return &res
}

func New(opts ...Opt) Parser {
	var res parser
	res.config = defaultConfig()
	for _, opt := range opts {
		opt(res.config)
	}
	return &res
}

func (jw *parser) newFilePath(relPath string, param map[string]interface{}) (string, string) {
	dirPath, fileName := filepath.Split(relPath)

	if !strings.Contains(fileName, ".tmpl") {
		// current only support tmpl files.
		return "", ""
	}

	newFileName := strings.ReplaceAll(fileName, ".tmpl", "")

	for k, v := range param {
		switch s := v.(type) {
		case string:
			newFileName = strings.Replace(newFileName, fmt.Sprintf("{{%s}}", k), s, 1)
			dirPath = strings.Replace(dirPath, fmt.Sprintf("{{%s}}", k), s, 1)
		}
	}

	newFileName = strings.ReplaceAll(newFileName, "dot@", ".")
	dirPath = strings.ReplaceAll(dirPath, "dot@", ".")

	return dirPath, newFileName
}

func (jw *parser) Parse(source, destination string) (err error) {
	var wg sync.WaitGroup
	var errch = make(chan error, 1)
	// discover the root first, the travelPath is at the root
	err = fs.WalkDir(jw.fs, source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var (
			param              = jw.params
			relPath            = strings.Replace(path, source, "", 1)
			destPath, destFile = jw.newFilePath(relPath, param)
		)
		if d.IsDir() {
			dirPath := strings.ReplaceAll(destPath, "dot@", ".")
			_ = os.MkdirAll(dirPath, os.ModePerm)
			return nil
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			//var param = jw.params.createFlexCase()
			var destFilePath = filepath.Join(destination, destPath, destFile)

			if err := os.MkdirAll(filepath.Join(destination, destPath), os.ModePerm); err != nil {
				errch <- err
			}
			if _, err := os.Stat(destFilePath); !os.IsNotExist(err) {
				if jw.force || (jw.override && strings.Contains(destFile, "_gen")) {
					log.Printf("Update %s\n", destFilePath)
				} else {
					errch <- err
					return
				}
			} else {
				log.Printf("Generate %s\n", destFilePath)
			}
			fo, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				errch <- err
				return
			}
			defer func() {
				_ = fo.Close()
			}()
			if len(jw.watermark) > 0 && filepath.Ext(destFile) == ".go" {
				_, _ = fo.WriteString(jw.watermark)
			}
			tmpl, err := template.ParseFS(jw.fs, path)
			if err != nil {
				errch <- err
				return
			}

			var buf = new(bytes.Buffer)
			if err := tmpl.Execute(buf, param); err != nil {
				errch <- err
				return
			}

			b := buf.Bytes()

			if filepath.Ext(destFile) == "go" {
				b, err = format.Source(b, format.Options{
					LangVersion: "1.17",
				})
				if err != nil {
					errch <- err
					return
				}
			}
			_, _ = fo.Write(b)
		}()
		return nil
	})

	go func() {
		wg.Wait()
		close(errch)
	}()
	for err := range errch {
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
	return err
}
