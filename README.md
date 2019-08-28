# Go Templater
A simple Golang Template module with auto reload and tree like file layout.

## Usage
By default Go Templater will read from "view" directory and will output it's build to "cache/view" but can be overriden in the config. You can also disable the Hot Reload feature in the config.

Here is the directory structure for the view : 
```
view/
    -home/
        -deep.html
    -part/
        -greetings.html
        -who.html
    -home.html
```
And here is the minimum working code to use Go Templater :
```
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
```

## Example
For more information please see [Example](example).

## Other Variant
There is also other variant of Go Templater that works with Gin Web Framework, you can find it [HERE](https://github.com/firmanmm/gin-templater)