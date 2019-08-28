package gotemplater

import (
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type event func(string)
type signal uint

var stopSignal = signal(0)

type templaterWatcher struct {
	logger     *log.Logger
	watchDir   string
	signalChan chan signal
	watcher    *fsnotify.Watcher
	onEvent    []event
	isRunning  bool
}

func (t *templaterWatcher) broadcastEvent(path string) {
	for _, ev := range t.onEvent {
		ev(path)
	}
}

func (t *templaterWatcher) work() {
	t.isRunning = true
	timeThreshold := time.Now().Add(time.Second)
	for {
		select {
		case _ = <-t.signalChan:
			t.isRunning = false
			close(t.signalChan)
			runtime.Goexit()
		case ev := <-t.watcher.Events:
			if ev.Op == fsnotify.Write {
				if time.Now().After(timeThreshold) {
					timeThreshold = time.Now().Add(time.Second)
					ev.Name = strings.ReplaceAll(ev.Name, "\\", "/")
					t.broadcastEvent(ev.Name)
				}
			} else if ev.Op == fsnotify.Create {
				t.watcher.Add(ev.Name)
			}
		}
	}
}

func (t *templaterWatcher) initDirectoryWatcher(path string) {
	var stack []string
	stack = append(stack, path)
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		t.watcher.Add(current)
		file, err := os.Open(current)
		if err != nil {
			log.Fatalln(err)
		}
		fileInfos, err := file.Readdir(0)
		file.Close()
		if err != nil {
			log.Fatalln(err)
		}
		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() {
				stack = append(stack, current+"/"+fileInfo.Name())
			}
		}
	}
}

func (t *templaterWatcher) Run() {
	if t.isRunning {
		return
	}
	t.signalChan = make(chan signal)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.logger.Fatalln(err)
	}
	t.watcher = watcher
	t.initDirectoryWatcher(t.watchDir)
	go t.work()
}

func (t *templaterWatcher) Stop() {
	if !t.isRunning {
		return
	}
	t.signalChan <- stopSignal
}

func (t *templaterWatcher) addListener(ev event) {
	t.onEvent = append(t.onEvent, ev)
}

func newTemplaterWatcher(dir string, logger *log.Logger) *templaterWatcher {
	instance := new(templaterWatcher)
	instance.watchDir = dir
	instance.logger = logger
	instance.onEvent = make([]event, 0)
	return instance
}
