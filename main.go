package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

const (
	ExitCodeOK    = 0
	ExitCodeError = 1
)

type CLI struct {
	outStream, errStream io.Writer
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

func (c *CLI) Run(args []string) int {
	files := make(map[string]gistContent)

	if len(args) > 1 {
		for _, fileName := range args[1:] {
			f, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintln(c.errStream, err)
				return ExitCodeError
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Fprintln(c.errStream, err)
				return ExitCodeError
			}

			files[fileName] = gistContent{Content: string(b)}
		}
	} else {
		// 標準入力を待つ
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(c.errStream, err)
			return ExitCodeError
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
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	resj := make(map[string]interface{})

	err = json.Unmarshal(body, &resj)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	fmt.Fprintln(c.outStream, resj["html_url"])

	return ExitCodeOK
}
