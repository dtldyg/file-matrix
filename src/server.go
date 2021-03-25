package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/reg", reg)
	http.HandleFunc("/index", index)
	log.Fatal(http.ListenAndServe("0.0.0.0:9989", nil))
}

func reg(resp http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(b))
	time.Sleep(time.Second)
	_, _ = resp.Write([]byte("renew"))
}

func index(resp http.ResponseWriter, req *http.Request) {
}

// --- utils ---
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
