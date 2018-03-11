package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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
	gist := &github.Gist{Files: make(map[github.GistFilename]github.GistFile)}

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

			bStr := string(b)
			gist.Files[github.GistFilename(fileName)] = github.GistFile{
				Filename: &fileName,
				Content:  &bStr,
			}
		}
	} else {
		// 標準入力を待つ
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(c.errStream, err)
			return ExitCodeError
		}

		bStr := string(b)
		gist.Files[github.GistFilename("")] = github.GistFile{
			Content: &bStr,
		}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	gi, _, err := client.Gists.Create(ctx, gist)

	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	fmt.Fprintln(c.outStream, *gi.HTMLURL)

	return ExitCodeOK
}

func getToken() string {
	token := os.Getenv("GITHUB_TOKEN")

	if token != "" {
		return token
	}

	// TOML

	return ""
}
