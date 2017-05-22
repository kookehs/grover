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
	dc.Crawl("C:/", options.verbose)

	for len(dc.frontier) > 0 {
		directory := dc.frontier[len(dc.frontier)-1]
		dc.frontier = dc.frontier[:len(dc.frontier)-1]
		dc.Crawl(directory, options.verbose)
	}

	fmt.Println("File count: ", len(dc.files))

	for _, file := range dc.Find(-1, options.target, 0.2) {
		fmt.Println(file.path + file.name)
	}
}
