package gotemplater

import (
	"html/template"
)

//Config hold configuration for Templater
type Config struct {
	AutoReload bool
	InputDir   string
	OutputDir  string
	FuncMap    template.FuncMap
}

//NewConfig return a new instance of Config with following properties
//AutoReload set to true
//InputDir set to view
//OutputDir set to cache/view
func NewConfig() *Config {
	instance := new(Config)
	instance.InputDir = "view"
	instance.OutputDir = "cache/view"
	instance.AutoReload = true
	instance.FuncMap = template.FuncMap{}
	return instance
}
