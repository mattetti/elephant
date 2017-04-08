package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"

	"flag"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattetti/elephant"
)

var (
	flagCachePath = flag.String("path", filepath.Join("cache", "goweekly"), "path to the folder where to store the golang weekly cached pages")
)

func main() {
	flag.Parse()
	for _, path := range cachedIssues() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("failed to open %s - %v\n", path, err)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(f)
		if err != nil {
			f.Close()
			fmt.Println("failed to read", path, err)
			continue
		}
		issue := elephant.Issue{Doc: doc}
		err = issue.Parse()
		f.Close()
		if err != nil {
			fmt.Println("failed to parse issue", path, err)
			continue
		}
		fmt.Printf("\n*** Issue %s ***\n\n", path)
		for _, item := range issue.Items {
			fmt.Printf("%+v\n", item)
		}
	}
}

func cachedIssues() []string {
	files, err := ioutil.ReadDir(*flagCachePath)
	if err != nil {
		log.Fatal(err)
	}

	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = filepath.Join(*flagCachePath, file.Name())
	}
	return paths
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
