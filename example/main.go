package main

import (
	"log"

	gotemplater "github.com/firmanmm/go-templater"
)

func main() {
	config := gotemplater.NewConfig()
	templater := gotemplater.NewTemplater(config)
	templater.Run()

	res, err := templater.RenderToString("home.html", map[string]interface{}{
		"message": "Hello!",
		"who":     "A Message",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(res)

	res, err = templater.RenderToString("home.deep.html", map[string]interface{}{
		"message": "Hello!",
		"who":     "A Deeper Message",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(res)
}
