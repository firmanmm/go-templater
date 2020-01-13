package gotemplater

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"sync"
)

//Templater instance to perform hot reload on Go instance
type Templater struct {
	logger     *log.Logger
	template   *template.Template
	builder    *templaterBuilder
	watcher    *templaterWatcher
	hotReload  bool
	outputDir  string
	bufferPool *sync.Pool
	funcMap    template.FuncMap
}

func (t *Templater) Run() {
	t.builder.initBuild()
	t.reload()
	if t.hotReload {
		t.watcher.Run()
		t.logger.Println("Templater started")
	}
}

func (t *Templater) Stop() {
	if t.hotReload {
		t.watcher.Stop()
		t.logger.Println("Templater stopped")
	}
}

func (t *Templater) reload() {
	newTemplate, err := template.New("").
		Funcs(t.funcMap).
		ParseGlob(t.outputDir + "/*")
	if err != nil {
		log.Println(err.Error())
		return
	}
	t.template = newTemplate
	t.logger.Println("Templater reloaded")
}

func (t *Templater) Render(wr io.Writer, name string, data interface{}) error {
	return t.template.ExecuteTemplate(wr, name, data)
}

func (t *Templater) RenderToByteArray(name string, data interface{}) ([]byte, error) {
	buffer := t.bufferPool.Get().(*bytes.Buffer)
	if err := t.Render(buffer, name, data); err != nil {
		return nil, err
	}
	defer buffer.Reset()
	resBuffer := make([]byte, buffer.Len())
	copy(resBuffer, buffer.Bytes())
	defer t.bufferPool.Put(buffer)
	return resBuffer, nil
}

func (t *Templater) RenderToString(name string, data interface{}) (string, error) {
	buffer := t.bufferPool.Get().(*bytes.Buffer)
	if err := t.Render(buffer, name, data); err != nil {
		return "", err
	}
	defer buffer.Reset()
	defer t.bufferPool.Put(buffer)
	return buffer.String(), nil
}

//NewTemplater will return a new and ready to use Templater instance
func NewTemplater(conf *Config) *Templater {
	instance := new(Templater)
	instance.outputDir = conf.OutputDir
	instance.hotReload = conf.AutoReload
	instance.logger = log.New(os.Stdout, "[GO-TEMPLATER] ", log.Ltime)
	instance.builder = newTemplaterBuilder(conf.InputDir, conf.OutputDir, instance.logger)
	instance.watcher = newTemplaterWatcher(conf.InputDir, instance.logger)
	instance.funcMap = conf.FuncMap
	instance.bufferPool = new(sync.Pool)
	instance.bufferPool.New = func() interface{} {
		byteBuffer := make([]byte, 0)
		return bytes.NewBuffer(byteBuffer)
	}
	if conf.AutoReload {
		rebuildEv := func(data string) {
			instance.reload()
		}
		instance.watcher.addListener(instance.builder.generate)
		instance.watcher.addListener(rebuildEv)
	}
	return instance
}
