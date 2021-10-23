package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vchitai/x-crafter/parser"
)

type Builder struct {
	*Guide
	parser parser.Parser
}

func (builder *Builder) getParser() parser.Parser {
	if builder.parser != nil {
		return builder.parser
	}

	var fsOpt = builder.TemplateFS
	if fsOpt == nil {
		fsOpt = os.DirFS(builder.TemplateRoot)
	}
	return parser.New(
		parser.WithWatermark(builder.Watermark),
		parser.WithFS(fsOpt),
		parser.WithParams(builder.Params),
	)
}
func (builder *Builder) execute(step *Step, at string) error {
	var startTime = time.Now()
	log.Println("Working on step", step.Name)
	at = filepath.Join(at, step.On)
	if step.Condition != nil {
		if step.Condition.NotInstalled != "" {
			if _, err := exec.LookPath(step.Condition.NotInstalled); err == nil {
				// Found
				return nil
			}
		}
		if step.Condition.Exists == "" {
			if x, err := filepath.Glob(step.Condition.Exists); err != nil {
				return nil
			} else if len(x) == 0 {
				// not doing
				return nil
			}
		}
	}
	if len(step.Run) > 0 {
		var args = step.Run[1:]
		var envM = envToMap(step.Env)
		for idx, arg := range args {
			args[idx] = os.Expand(arg, func(s string) string {
				return envM[s]
			})
		}
		cmd := newCmd(step.Run[0], step.Run[1:]...)
		if len(step.Env) > 0 {
			cmd.Env = append(os.Environ(), step.Env...)
		}
		cmd.Dir = at
		if os.Getenv("DEBUG") != "false" {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		return cmd.Run()
	}
	if step.Parse != "" {
		if step.Repeat != nil {
			for _, m := range interfaceToMapSlice(builder.Params[step.Repeat.For]) {
				var subParam = cloneMap(builder.Params)
				for k, v := range m {
					subParam[k] = v
				}
				var psr = builder.getParser().WithOption(
					parser.WithParams(subParam),
				)
				if err := psr.Parse(step.Parse, at); err != nil {
					log.Println("Parsing error", err)
				}
			}
		} else {
			var psr = builder.getParser()
			if err := psr.Parse(step.Parse, at); err != nil {
				log.Println("Parsing error", err)
			}
		}
	}
	log.Printf("Done in %0.3f\n", time.Since(startTime).Seconds())
	return nil
}

func (builder *Builder) Execute(at string) error {
	if err := os.MkdirAll(at, os.ModePerm); err != nil {
		return err
	}

	var env []string
	for k, v := range builder.Params {
		switch v.(type) {
		case string:
			env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
		}
	}
	for _, step := range builder.Steps {
		step.Env = append(step.Env, env...)
		if err := builder.execute(step, at); err != nil {
			return err
		}
	}
	return nil
}
func New(guide *Guide) *Builder {
	return &Builder{Guide: guide}
}
