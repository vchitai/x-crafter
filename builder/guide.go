package builder

import (
	"io/fs"
	"strings"
)

type StepCondition struct {
	NotInstalled string `json:"not_installed" yaml:"not_installed"`
	Exists       string `json:"exists" yaml:"exists"`
}
type StepRepeat struct {
	For string `json:"for" yaml:"for"`
}
type Step struct {
	Run       []string       `json:"run" yaml:"run,flow"` // Cmd for run
	Env       []string       `json:"env" yaml:"env"`
	Parse     string         `json:"parse" yaml:"parse"`
	On        string         `json:"on" yaml:"on"`
	Name      string         `json:"name" yaml:"name"`
	Condition *StepCondition `json:"condition" yaml:"condition"`
	Repeat    *StepRepeat    `json:"repeat" yaml:"repeat"`
}

type Guide struct {
	Steps        []*Step                `json:"steps" yaml:"steps"`
	Params       map[string]interface{} `json:"params" yaml:"params"`
	TemplateRoot string                 `json:"template_root" yaml:"template_root"`
	Watermark    string                 `json:"watermark" yaml:"watermark"`
	TemplateFS   fs.FS                  `json:"-" yaml:"-"`
}

func defaultGuide() *Guide {
	return &Guide{TemplateRoot: "./layers"}
}

func envToMap(env []string) map[string]string {
	var res = make(map[string]string)
	for _, e := range env {
		ss := strings.Split(e, "=")
		if len(ss) != 2 {
			continue
		}
		res[ss[0]] = ss[1]
	}
	return res
}
