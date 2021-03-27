package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	serverAddr = "http://127.0.0.1:9989/"
)

var (
	name = ""
	key  = ""
)

func main() {
	name, key = nameKey()
	data := map[string][]string{
		"name": {name},
	}
	for {
		resp, err := http.PostForm(serverAddr+"reg", data)
		if err != nil {
			panic(err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
	}
}

func nameKey() (string, string) {
	name, key := "", ""
	fmt.Print("name:")
	_, _ = fmt.Scan(&name)
	fmt.Print("key:")
	_, _ = fmt.Scan(&key)
	return name, key
}

func reg() {
}

func push(file) {
}
