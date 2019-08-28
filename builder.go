package gotemplater

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type templaterBuilder struct {
	logger    *log.Logger
	inputDir  string
	outputDir string
}

func (b *templaterBuilder) initBuild() {
	b.prepareOutputDir()

	var stack []string
	stack = append(stack, b.inputDir)
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		file, err := os.Open(current)
		if err != nil {
			b.logger.Fatalln(err)
		}
		defer file.Close()
		fileInfos, err := file.Readdir(0)
		if err != nil {
			b.logger.Fatalln(err)
		}
		for _, fileInfo := range fileInfos {
			currentPath := fmt.Sprintf("%s/%s", current, fileInfo.Name())
			if fileInfo.IsDir() {
				stack = append(stack, currentPath)
			} else {
				b.generate(currentPath)
			}
		}
	}
}

func (b *templaterBuilder) prepareOutputDir() {
	if err := os.MkdirAll(b.outputDir, 0666); err != nil {
		b.logger.Fatalln(err)
	}
}

func (b *templaterBuilder) generate(path string) {
	buildPath := path[len(b.inputDir)+1:]
	buildName := strings.ReplaceAll(buildPath, "/", ".")
	linkName := fmt.Sprintf("%s/%s", b.outputDir, buildName)
	os.Link(path, linkName)
}

func newTemplaterBuilder(inputDir, outputDir string, logger *log.Logger) *templaterBuilder {
	instance := new(templaterBuilder)
	instance.inputDir = inputDir
	instance.outputDir = outputDir
	instance.logger = logger
	return instance
}
