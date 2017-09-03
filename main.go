package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type gistContent struct {
	Content string `json:"content"`
}

type gist struct {
	Description string                 `json:"description"`
	Public      bool                   `json:"public"`
	Files       map[string]gistContent `json:"files"`
}

func main() {
	var content string
	var fileName string

	if len(os.Args) > 1 {
		fileName = os.Args[1]
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		content = string(b)
	} else {
		// 標準入力を待つ
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		content = string(b)
	}

	jsonObj := &gist{
		Description: "",
		Public:      false,
		Files:       map[string]gistContent{fileName: gistContent{Content: content}},
	}

	b, _ := json.Marshal(jsonObj)

	res, err := http.Post("https://api.github.com/gists", "application/json;charset=utf-8", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		log.Fatal(string(body))
	}

	resj := make(map[string]interface{})

	err = json.Unmarshal(body, &resj)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resj["html_url"])
}
