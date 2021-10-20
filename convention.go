package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Convention struct {
	Arguments []string `yaml:"arguments"`
	Basket    struct {
		Map map[string][]string `yaml:"map"`
	} `yaml:"basket"`
	Syntax map[string]struct {
		Start string `yaml:"start"`
		End   string `yaml:"end"`
	} `yaml:"syntax"`
}

func makeReplacementSet(convention *Convention) *strings.Replacer {
	if convention == nil {
		return nil
	}

	var replacementSets = make([]string, 0)

	replacementSets = append(replacementSets, "/* end */", "{{ end }}")

	if len(convention.Arguments) > 0 {
		for _, arg := range convention.Arguments {
			replacementSets = append(replacementSets, arg, fmt.Sprintf("{{.%s}}", arg))
		}
	}

	return strings.NewReplacer(replacementSets...)
}

func makeBasket(convention *Convention) map[string]string {
	if convention == nil {
		return make(map[string]string)
	}
	var res = make(map[string]string)
	for layer, globs := range convention.Basket.Map {
		for _, glob := range globs {
			res[glob] = layer
		}
	}
	return res
}

func makeSyntaxRegexp(convention *Convention) map[string]*regexp.Regexp {
	if convention == nil {
		return make(map[string]*regexp.Regexp)
	}
	var res = make(map[string]*regexp.Regexp, len(convention.Syntax))
	for ext, rule := range convention.Syntax {
		res[ext] = regexp.MustCompile(regexp.QuoteMeta(rule.Start) + "(.*)" + regexp.QuoteMeta(rule.End))
	}
	return res
}

func loadConvention(source string) *Convention {
	f, err := os.Open(filepath.Join(source, "convention.yaml"))
	if err != nil {
		log.Println("Cannot load convention", err)
		return nil
	}
	defer func() {
		_ = f.Close()
	}()

	rawConvention, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("Cannot load convention", err)
		return nil
	}
	var convention Convention
	if err := yaml.Unmarshal(rawConvention, &convention); err != nil {
		log.Println("Cannot load convention", err)
		return nil
	}

	return &convention
}
