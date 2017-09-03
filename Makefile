.PHONY: release

all: gist

gist: main.go
	go build -o gist main.go

release:
	GOOS=linux go build -o gist main.go
	tar cvf gist-linux-amd64.tar.gz gist
	GOOS=darwin go build -o gist main.go
	tar cvf gist-darwin-amd64.tar.gz gist
