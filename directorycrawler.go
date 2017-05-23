package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

// Match contains information regarding a potential target
type Match struct {
	file       *File
	confidence float64
}

// Matches is a type alias for an array of pointers of type Match
type Matches []*Match

// NewMatch returns a pointer to an initialized Match
func NewMatch() *Match {
	m := new(Match)
	m.file = nil
	m.confidence = 0
	return m
}

// Len returns the length of the slice
func (slice Matches) Len() int {
	return len(slice)
}

// Less returns whether the element at i is less than the element at j
func (slice Matches) Less(i, j int) bool {
	return slice[i].confidence < slice[j].confidence
}

// Swap switches around the element at i with the element at j
func (slice Matches) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

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
	files           []*File
	frontier        []string
	ignoreDirectory []string
	ignoreExtension []string
}

// NewDirectoryCrawler returns a pointer to an initialized DirectoryCrawler
func NewDirectoryCrawler() *DirectoryCrawler {
	dc := new(DirectoryCrawler)
	dc.files = make([]*File, 0)
	dc.frontier = make([]string, 0)
	dc.ignoreDirectory = make([]string, 0)
	dc.ignoreExtension = make([]string, 0)
	return dc
}

// CacheResults stores results on disk for faster searches
func (dc *DirectoryCrawler) CacheResults(directory string, verbose bool) {
	home := os.Getenv("HOME")

	if home != "" {
		err := os.MkdirAll(home+"/.cache/grover/", os.ModeDir)

		if err != nil {
			Println("[Error] "+err.Error(), verbose)
		} else {
			sanitized := dc.SanitizePath(directory)
			file, err := os.Create(home + "/.cache/grover/" + sanitized + ".txt")

			if err != nil {
				Println("[Error] "+err.Error(), verbose)
			} else {
				writer := bufio.NewWriter(file)

				for _, v := range dc.files {
					_, err := writer.WriteString(v.path + v.name + "\n")

					if err != nil {
						Println("[Error] "+err.Error(), verbose)
					}
				}

				writer.Flush()
				file.Close()
			}
		}
	}
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

		for _, v := range dc.ignoreDirectory {
			if path.Base(v) == path.Base(directory) {
				ignore = true
				break
			}
		}

		if !ignore {
			Println("[Crawling] "+directory, verbose)
			items, err := ioutil.ReadDir(directory)

			if err != nil {
				Println("[Error] "+err.Error(), verbose)
			} else {
				dc.Visit(directory, items, verbose)
			}
		} else {
			Println("[Ignoring] "+directory, verbose)
		}
	}

	dc.CacheResults(directory, verbose)
}

// Find returns an array of possible matches
func (dc *DirectoryCrawler) Find(limit int, name string, threshold float64) Matches {
	matches := make(Matches, 0)

	for _, file := range dc.files {
		confidence := FuzzySearch(file.name, name)

		if confidence > threshold {
			match := NewMatch()
			match.file = file
			match.confidence = confidence
			matches = append(matches, match)
		}

		if limit != -1 && len(matches) == limit {
			break
		}
	}

	sort.Sort(matches)
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
			line := scanner.Text()

			if line[0] == '*' {
				dc.ignoreExtension = append(dc.ignoreExtension, line)
			} else {
				dc.ignoreDirectory = append(dc.ignoreDirectory, line)
			}
		}

		file.Close()
	}
}

// SanitizePath replaces all invalid filename characters with percents
func (dc *DirectoryCrawler) SanitizePath(path string) string {
	path = strings.Replace(path, "\\", "%", -1)
	path = strings.Replace(path, "/", "%", -1)
	path = strings.Replace(path, ":", "%", -1)
	path = strings.Replace(path, "*", "%", -1)
	path = strings.Replace(path, "?", "%", -1)
	path = strings.Replace(path, "\"", "%", -1)
	path = strings.Replace(path, "<", "%", -1)
	path = strings.Replace(path, ">", "%", -1)
	path = strings.Replace(path, "|", "%", -1)
	return path
}

// Visit tracks files and folders in the specified directory
func (dc *DirectoryCrawler) Visit(directory string, items []os.FileInfo, verbose bool) {
	for _, item := range items {
		switch item.IsDir() {
		case false:
			ignore := false

			for _, v := range dc.ignoreExtension {
				if path.Ext(v) == path.Ext(item.Name()) {
					ignore = true
					break
				}
			}

			if !ignore {
				file := NewFile()
				file.name = item.Name()
				file.path = directory
				dc.files = append(dc.files, file)
			} else {
				Println("[Ignoring] "+directory+item.Name(), verbose)
			}
		case true:
			dc.frontier = append(dc.frontier, directory+item.Name())
		}
	}
}
