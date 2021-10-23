package parser

import (
	"embed"
	"io/fs"
)

type config struct {
	fs        fs.FS
	params    map[string]interface{}
	override  bool
	force     bool
	watermark string
}
type Opt func(c *config)

func WithFS(fs fs.FS) Opt {
	return func(c *config) {
		c.fs = fs
	}
}

func WithParams(params map[string]interface{}) Opt {
	return func(c *config) {
		c.params = params
	}
}
func WithOverride(override bool) Opt {
	return func(c *config) {
		c.override = override
	}
}

func WithWatermark(watermark string) Opt {
	return func(c *config) {
		c.watermark = watermark
	}
}

func defaultConfig() *config {
	return &config{
		fs:       embed.FS{},
		params:   nil,
		override: false,
		force:    false,
	}
}
