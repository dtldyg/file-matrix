package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	serverAddr = "http://127.0.0.1:9989/"
)

func main() {
	fmt.Print("name:")
	name := ""
	_, _ = fmt.Scan(&name)
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
