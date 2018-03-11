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
	os.Exit(run(os.Args))
}

func run(args []string) int {
	files := make(map[string]gistContent)

	if len(args) > 1 {
		for _, fileName := range args[1:] {
			f, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}

			files[fileName] = gistContent{Content: string(b)}
		}
	} else {
		// 標準入力を待つ
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		files[""] = gistContent{Content: string(b)}
	}

	jsonObj := &gist{
		Description: "",
		Public:      false,
		Files:       files,
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

	return 0
}
