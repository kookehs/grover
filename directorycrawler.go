package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
)

// File contains location information regarding the file
type File struct {
	name string
	path string
}

// NewFile returns a pointer to an initialized File
func NewFile() *File {
	file := new(File)
	file.name = ""
	file.path = ""
	return file
}

// DirectoryCrawler contains files regarding the crawled directory
type DirectoryCrawler struct {
	files    []*File
	frontier []string
	ignore   []string
}

// NewDirectoryCrawler returns a pointer to an initialized DirectoryCrawler
func NewDirectoryCrawler() *DirectoryCrawler {
	dc := new(DirectoryCrawler)
	dc.files = make([]*File, 0)
	dc.frontier = make([]string, 0)
	dc.ignore = make([]string, 0)
	return dc
}

// Crawl walks through each directory calling Visit each time
func (dc *DirectoryCrawler) Crawl(directory string, verbose bool) {
	dc.frontier = append(dc.frontier, directory)

	for len(dc.frontier) > 0 {
		directory := dc.frontier[len(dc.frontier)-1]
		dc.frontier = dc.frontier[:len(dc.frontier)-1]

		if directory[len(directory)-1] != '/' {
			directory += "/"
		}

		ignore := false

		for _, v := range dc.ignore {
			if path.Base(v) == path.Base(directory) {
				ignore = true
				break
			}
		}

		if !ignore {
			Println("Crawling: "+directory, verbose)
			items, err := ioutil.ReadDir(directory)

			if err != nil {
				Println(err, verbose)
			} else {
				dc.Visit(directory, items)
			}
		} else {
			Println("Ignoring: "+directory, verbose)
		}
	}
}

// Find returns an array of possible matches
func (dc *DirectoryCrawler) Find(limit int, name string, threshold float64) []*File {
	matches := make([]*File, 0)

	for _, file := range dc.files {
		if FuzzySearch(name, file.path) > threshold || FuzzySearch(name, file.name) > threshold {
			matches = append(matches, file)
		}

		if limit != -1 && len(matches) == limit {
			break
		}
	}

	return matches
}

// LoadIgnore reads the .groverignore in the current directory
func (dc *DirectoryCrawler) LoadIgnore(verbose bool) {
	file, err := os.Open("./.groverignore")

	if err != nil {
		Println(err, verbose)
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			dc.ignore = append(dc.ignore, scanner.Text())
		}
	}

	file.Close()
}

// Visit tracks files and folders in the specified directory
func (dc *DirectoryCrawler) Visit(directory string, items []os.FileInfo) {
	for _, item := range items {
		switch item.IsDir() {
		case false:
			file := NewFile()
			file.name = item.Name()
			file.path = directory
			dc.files = append(dc.files, file)
		case true:
			dc.frontier = append(dc.frontier, directory+item.Name())
		}
	}
}
