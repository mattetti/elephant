.DEFAULT_GOAL := default

goweeklydl: cmd/dlgoweekly/main.go *.go
	go build -o goweeklydl github.com/mattetti/elephant/cmd/dlgoweekly

goweeklyprocess: cmd/processgoweekly/main.go *.go
	go build -o goweeklyprocess github.com/mattetti/elephant/cmd/processgoweekly

default: goweeklydl goweeklyprocess
	./goweeklydl
	./goweeklyprocess