package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// client
	http.HandleFunc("/reg", reg)
	http.HandleFunc("/push", push)
	// browser
	http.HandleFunc("/index", index)
	http.HandleFunc("/names", names)
	http.HandleFunc("/dir", dir)
	http.HandleFunc("/file", file)
	log.Fatal(http.ListenAndServe("0.0.0.0:9989", nil))
}

// client
func reg(resp http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(b))
	time.Sleep(time.Second)
	_, _ = resp.Write([]byte("renew"))
}

func push(resp http.ResponseWriter, req *http.Request) {
}

// browser
func index(resp http.ResponseWriter, req *http.Request) {
}

func names(resp http.ResponseWriter, req *http.Request) {
}

func dir(resp http.ResponseWriter, req *http.Request) {
}

func file(resp http.ResponseWriter, req *http.Request) {
}

// ----- utils -----
func sizeToString(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%dB", n)
	} else if n < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(n)/float64(1024))
	} else if n < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(n)/float64(1024*1024))
	} else {
		return fmt.Sprintf("%.2fGB", float64(n)/float64(1024*1024*1024))
	}
}
