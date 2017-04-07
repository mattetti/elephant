package main

import "github.com/mattetti/elephant"

func main() {
	issue := elephant.Issue{Nbr: 41}
	if err := issue.Fetch(); err != nil {
		panic(err)
	}
}
