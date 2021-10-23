package breaker

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

type Guide struct {
	KwArgs map[string]string   `yaml:"kwargs"`
	Args   []string            `yaml:"args"`
	Layers map[string][]string `yaml:"layers"`
	Syntax map[string]struct {
		Start    string `yaml:"start"`
		End      string `yaml:"end"`
		Inline   string `yaml:"inline"`
		Optional string `yaml:"optional"`
	} `yaml:"syntax"`
}

func makeContentReplacementSet(convention *Guide) *strings.Replacer {
	if convention == nil {
		return nil
	}

	// TODO: short by len desc
	var replacementSets = make([]string, 0)
	if len(convention.Args) > 0 {
		for _, arg := range convention.Args {
			replacementSets = append(replacementSets, arg, fmt.Sprintf("{{.%s}}", arg))
		}
	}
	if len(convention.KwArgs) > 0 {
		for kwarg, replacement := range convention.KwArgs {
			replacementSets = append(replacementSets, kwarg, fmt.Sprintf("{{.%s}}", replacement))
		}
	}

	return strings.NewReplacer(replacementSets...)
}

func makeFilenameReplacementSet(convention *Guide) *strings.Replacer {
	if convention == nil {
		return nil
	}

	// TODO: short by len desc
	var replacementSets = make([]string, 0)
	if len(convention.Args) > 0 {
		for _, arg := range convention.Args {
			replacementSets = append(replacementSets, arg, fmt.Sprintf("{{%s}}", arg))
		}
	}
	if len(convention.KwArgs) > 0 {
		for kwarg, replacement := range convention.KwArgs {
			replacementSets = append(replacementSets, kwarg, fmt.Sprintf("{{%s}}", replacement))
		}
	}

	return strings.NewReplacer(replacementSets...)
}

func makeLayerIndex(convention *Guide) map[string][]string {
	if convention == nil {
		return make(map[string][]string)
	}
	return convention.Layers
}

type matchRule struct {
	*regexp.Regexp
	newline    bool
	eof        bool
	needEscape bool
}

func makeSyntaxRegexp(convention *Guide) map[string][]*matchRule {
	if convention == nil {
		return make(map[string][]*matchRule)
	}
	var res = make(map[string][]*matchRule, len(convention.Syntax))
	for ext, rule := range convention.Syntax {
		if len(rule.Start) > 0 && len(rule.End) > 0 {
			res[ext] = append(res[ext], &matchRule{
				Regexp:  regexp.MustCompile(regexp.QuoteMeta(rule.Start) + "(.*)" + regexp.QuoteMeta(rule.End)),
				newline: false,
			})
		}
		if len(rule.Inline) > 0 {
			res[ext] = append(res[ext], &matchRule{
				Regexp:  regexp.MustCompile(regexp.QuoteMeta(rule.Inline) + "(.*)\\r*\\n"),
				newline: true,
			})
			res[ext] = append(res[ext], &matchRule{
				Regexp: regexp.MustCompile(regexp.QuoteMeta(rule.Inline) + "(.*)$"),
				eof:    true,
			})
		}
		if len(rule.Optional) > 0 {
			res[ext] = append(res[ext], &matchRule{
				Regexp:     regexp.MustCompile(regexp.QuoteMeta(rule.Optional) + "(.*)\\r*\\n"),
				newline:    true,
				needEscape: true,
			})
			res[ext] = append(res[ext], &matchRule{
				Regexp:     regexp.MustCompile(regexp.QuoteMeta(rule.Optional) + "(.*)$"),
				eof:        true,
				needEscape: true,
			})
		}
	}
	return res
}

func loadGuide(source string) *Guide {
	f, err := os.Open(filepath.Join(source, "xbreak.yaml"))
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
	var convention Guide
	if err := yaml.Unmarshal(rawConvention, &convention); err != nil {
		log.Println("Cannot load convention", err)
		return nil
	}
	return &convention
}
