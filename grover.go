package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args[1:]
	options := NewOptions()
	options.Parse(arguments)

	dc := NewDirectoryCrawler()
	dc.LoadIgnore(options.verbose)
	dc.Crawl("C:/", options)
	fmt.Println("File count: ", len(dc.files))

	for _, match := range dc.Find(-1, options.target, 0.6) {
		fmt.Println(match.file.path + match.file.name)
	}
}
