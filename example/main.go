package main

import (
	"html/template"
	"log"

	gotemplater "github.com/firmanmm/go-templater"
)

func main() {
	config := gotemplater.NewConfig()
	config.FuncMap = template.FuncMap{
		"hello": func(data string) string {
			return "Hello, " + data
		},
	}
	templater := gotemplater.NewTemplater(config)
	templater.Run()

	res, err := templater.RenderToString("home.html", map[string]interface{}{
		"message": "Hi There!",
		"who":     "A Message",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(res)

	res, err = templater.RenderToString("home.deep.html", map[string]interface{}{
		"message": "Deeper Hi There!",
		"who":     "A Deeper Message",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(res)
}
