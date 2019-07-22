package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Open returns the file contents as a string
func Open(filepath string) string {
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	return string(dat)
}

// Load returns the file contents as a string
func Load(filepath string) string {
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	return string(dat)
}

// Write writes the input string to a file
func Write(data, filepath string) {
	message := []byte(data)
	err := ioutil.WriteFile(filepath, message, 0644)
	check(err)
}

// Save writes the input string to a file
func Save(data, filepath string) {
	message := []byte(data)
	err := ioutil.WriteFile(filepath, message, 0644)
	check(err)
}

// List returns list of files in directory
func List(dir string) []string {
	files, err := ioutil.ReadDir(".")
	check(err)
	out := []string{}
	for _, file := range files {
		if !file.IsDir() {
			out = append(out, file.Name())
		}
	}
	return out
}

// ListAll Same as List() but scans subfolders too
func ListAll(dir string) []string {
	out := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		out = append(out, path)
		return nil
	})
	return out
}

// ListDir returns list of subdirectories in directory
func ListDir(dir string) []string {
	files, err := ioutil.ReadDir(".")
	check(err)
	out := []string{}
	for _, file := range files {
		if file.IsDir() {
			out = append(out, file.Name())
		}
	}
	return out
}

// SecureSave saves files and returns an event when write is confirmed
func SecureSave(data, filepath string, done chan bool) {
	Save(data, filepath)
	ndata := Open(filepath)
	if ndata == data {
		done <- true
	} else {
		SecureSave(data, filepath, done)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
