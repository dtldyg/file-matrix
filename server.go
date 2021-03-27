package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	nameStub map[string]*stub
	stubLock sync.RWMutex
	namePush map[string]*push
	pushLock sync.RWMutex
)

type stub struct {
	key     string
	optChan chan opt
}

type push struct {
	dirChan  chan []string
	fileChan chan []byte
}

type opt struct {
	cmd   int
	value string
}

func main() {
	nameStub = make(map[string]*stub)
	namePush = make(map[string]*push)
	// static
	http.Handle("/", http.FileServer(http.Dir(".")))
	// client
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/pushDir", pushDirHandler)
	http.HandleFunc("/pushFile", pushFileHandler)
	// browser
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/names", namesHandler)
	http.HandleFunc("/dir", dirHandler)
	http.HandleFunc("/file", fileHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:9989", nil))
}

// client
func regHandler(resp http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	key := req.FormValue("key")
	stub := &stub{
		key:     key,
		optChan: make(chan opt),
	}

	stubLock.Lock()
	nameStub[name] = stub
	stubLock.Unlock()

	select {
	case <-time.After(time.Second * 10):
		resp.Header().Set("opt", "renew")
	case opt := <-stub.optChan:
		switch opt.cmd {
		case 0:
			resp.Header().Set("opt", "dir")
		case 1:
			resp.Header().Set("opt", "file")
		}
		resp.Header().Set("value", opt.value)
	}

	stubLock.Lock()
	delete(nameStub, name)
	stubLock.Unlock()
}

func pushDirHandler(resp http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	path := req.FormValue("path")
	files := req.Form["dir"]

	pushLock.RLock()
	push := namePush[name]
	pushLock.RUnlock()

	push.dirChan <- append([]string{path}, files...)
}

func pushFileHandler(resp http.ResponseWriter, req *http.Request) {
	name := req.Header.Get("name")
	fileData, _ := ioutil.ReadAll(req.Body)

	pushLock.RLock()
	push := namePush[name]
	pushLock.RUnlock()

	push.fileChan <- fileData
}

// browser
func indexHandler(resp http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadFile("./index.html")
	_, _ = resp.Write(data)
}

func namesHandler(resp http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	key := req.FormValue("key")
	http.SetCookie(resp, &http.Cookie{Name: "user", Value: user})
	http.SetCookie(resp, &http.Cookie{Name: "key", Value: key})

	var names []string

	stubLock.RLock()
	for name, stub := range nameStub {
		if stub.key == key {
			names = append(names, name)
		}
	}
	stubLock.RUnlock()

	sort.Strings(names)
	namesHtml := ""
	for _, name := range names {
		namesHtml += fmt.Sprintf("<tr><td><a href=\"dir?name=%s&dir=\">%s</a></td></tr>\n", name, name)
	}
	data, _ := ioutil.ReadFile("./names.html")
	data = []byte(strings.ReplaceAll(string(data), "{names}", namesHtml))
	_, _ = resp.Write(data)
}

func dirHandler(resp http.ResponseWriter, req *http.Request) {
	userCoo, _ := req.Cookie("user")
	user := userCoo.Value
	keyCoo, _ := req.Cookie("key")
	key := keyCoo.Value
	name := req.FormValue("name")
	dir := req.FormValue("dir")

	stubLock.RLock()
	stub := nameStub[name]
	stubLock.RUnlock()

	if stub.key != key {
		return
	}
	_ = user

	push := &push{
		dirChan: make(chan []string),
	}
	pushLock.Lock()
	namePush[name] = push
	pushLock.Unlock()

	stub.optChan <- opt{
		cmd:   0,
		value: dir,
	}

	files := <-push.dirChan
	pushLock.Lock()
	delete(namePush, name)
	pushLock.Unlock()

	allDir := files[0]
	files = files[1:]
	filesHtml := ""
	if dir != "" {
		upDir := filepath.Dir(dir)
		if upDir == "." {
			upDir = ""
		}
		filesHtml += fmt.Sprintf("<tr><td><a href=\"dir?name=%s&dir=%s\">..</a></td></tr>\n", name, upDir)
	}
	for _, file := range files {
		isDir := file[0] == '0'
		rtName := file[1:]
		nameMod := strings.Split(rtName, "|")
		rtName = nameMod[0]
		rtMod := nameMod[1]
		rtSize := nameMod[2]
		rtPath := filepath.Join(dir, rtName)
		if isDir {
			filesHtml += fmt.Sprintf("<tr><td><a href=\"dir?name=%s&dir=%s\"><span class=\"glyphicon glyphicon-folder-open\" aria-hidden=\"true\"></span>\t%s</a></td><td>%s</td><td>%s</td></tr>\n", name, rtPath, rtName, rtMod, "-")
		} else {
			filesHtml += fmt.Sprintf("<tr><td><a href=\"file?name=%s&file=%s\"><span class=\"glyphicon glyphicon-download-alt\" aria-hidden=\"true\"></span>\t%s</a></td><td>%s</td><td>%s</td></tr>\n", name, rtPath, rtName, rtMod, rtSize)
		}
	}
	data, _ := ioutil.ReadFile("./dir.html")
	data = []byte(strings.ReplaceAll(string(data), "{dir}", allDir))
	data = []byte(strings.ReplaceAll(string(data), "{files}", filesHtml))
	_, _ = resp.Write(data)
}

func fileHandler(resp http.ResponseWriter, req *http.Request) {
	userCoo, _ := req.Cookie("user")
	user := userCoo.Value
	keyCoo, _ := req.Cookie("key")
	key := keyCoo.Value
	name := req.FormValue("name")
	file := req.FormValue("file")

	stubLock.RLock()
	stub := nameStub[name]
	stubLock.RUnlock()

	if stub.key != key {
		return
	}
	_ = user

	push := &push{
		fileChan: make(chan []byte),
	}
	pushLock.Lock()
	namePush[name] = push
	pushLock.Unlock()

	stub.optChan <- opt{
		cmd:   1,
		value: file,
	}

	fileData := <-push.fileChan
	pushLock.Lock()
	delete(namePush, name)
	pushLock.Unlock()

	_, fileName := filepath.Split(file)
	resp.Header().Set("Content-type", "application/octet-stream")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	_, _ = resp.Write(fileData)
}
