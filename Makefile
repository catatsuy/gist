all: gist

gist: main.go
	go build -o gist main.go
