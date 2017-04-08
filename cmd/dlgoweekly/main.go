// This CLI tool downloads all available goland weekly issues that
// aren't already available in the cache.
// The location of the cache can be set via the optional -path argument.
package main

import (
	"fmt"
	"io"
	"log"

	"net/http"
	"os"
	"path/filepath"

	"flag"

	"github.com/mattetti/elephant"
)

var (
	flagCachePath = flag.String("path", filepath.Join("cache", "goweekly"), "path to the folder where to store the golang weekly cached pages")
)

func main() {
	flag.Parse()
	fetchAllGoWeekly(*flagCachePath)
}

func fetchAllGoWeekly(cachePath string) {
	issueNbr := elephant.FirstGoWeeklyIssue
	var err error
	var resp *http.Response
	os.MkdirAll(cachePath, os.ModePerm)
	for err == nil {
		issuePath := filepath.Join(cachePath, fmt.Sprintf("issue_%d.html", issueNbr))
		if pathExists(issuePath) {
			fmt.Printf("skipped issue %d (already on disk)\n", issueNbr)
			issueNbr++
			continue
		}
		resp, err = http.Get(fmt.Sprintf(elephant.GoWeeklyIssuePage, issueNbr))
		if err != nil {
			break
		}
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("issue %d failed to download, status code: %d", issueNbr, resp.StatusCode)
			break
		}
		f, err := os.Create(issuePath)
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(f, resp.Body)
		resp.Body.Close()
		f.Close()
		if n < 42 {
			fmt.Println("Deleting empty issue", issueNbr)
			os.Remove(issuePath)
		} else {
			fmt.Println("downloaded issue", issueNbr)
		}
		issueNbr++
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
