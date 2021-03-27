package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	//serverAddr = "http://127.0.0.1:9989/"
	serverAddr = "http://149.129.82.53:9989/"
	rMod       = uint32(0b100100100)
)

var (
	root string
	name string
	key  string
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			select {}
		}
	}()
	root, _ = os.Getwd()
	name, key = initNameKey()
	for {
		opt, value := regReq()
		switch opt {
		case "renew":
			continue
		case "dir":
			dirReq(value)
		case "file":
			fileReq(value)
		}
	}
}

func initNameKey() (string, string) {
	name, key := "", ""
	fmt.Print("name:")
	_, _ = fmt.Scan(&name)
	fmt.Print("key:")
	_, _ = fmt.Scan(&key)
	return name, key
}

func regReq() (string, string) {
	resp, err := http.PostForm(serverAddr+"reg", url.Values{
		"name": {name},
		"key":  {key},
	})
	if err != nil {
		panic(err)
	}
	return resp.Header.Get("opt"), resp.Header.Get("value")
}

func dirReq(dirPath string) {
	curDir := filepath.Join(root, dirPath)
	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		panic(err)
	}

	var dir []string
	for _, file := range files {
		if !checkCanRead(file) || checkIsHidden(file) {
			continue
		}
		isDir := 1
		if file.IsDir() {
			isDir = 0
		}
		dir = append(dir, fmt.Sprintf("%d%s|%s|%s", isDir, file.Name(), file.ModTime().Format("2006-01-02 15:04:05"), sizeToString(file.Size())))
	}
	sort.Slice(dir, func(i, j int) bool { return strings.ToLower(dir[i]) < strings.ToLower(dir[j]) })
	_, err = http.PostForm(serverAddr+"pushDir", url.Values{
		"name": {name},
		"path": {curDir},
		"dir":  dir,
	})
	if err != nil {
		panic(err)
	}
}

func fileReq(filePath string) {
	curFile := filepath.Join(root, filePath)
	fileData, err := ioutil.ReadFile(curFile)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", serverAddr+"pushFile", bytes.NewReader(fileData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("name", name)
	_, _ = http.DefaultClient.Do(req)
}

// ----- utils -----
func sizeToString(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	} else if n < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(n)/float64(1024))
	} else if n < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(n)/float64(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(n)/float64(1024*1024*1024))
	}
}

func checkCanRead(file os.FileInfo) bool {
	return uint32(file.Mode())&rMod == rMod
}

func checkIsHidden(file os.FileInfo) bool {
	fa := reflect.ValueOf(file.Sys()).Elem().FieldByName("FileAttributes").Uint()
	b := []byte(strconv.FormatUint(fa, 2))
	if b[len(b)-2] == '1' {
		return true
	}
	return false
}
