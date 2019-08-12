package files

/* TODO:

Enable using URLs as either open or save targets,
Open would basically be a wrapper for an http.Get
Save would need token access to a cloud filesystem
service like Cloudvar

regex match the beginning to see if it is a web resource
if it is, http.Get, encode body as string, return

files.Login(token) // logs into cloudvar
if token is set
match the beginning of the url for cloudvar://
if it is, do an http transfer to the cloudvar servers
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Open returns the file contents as a string
func Open(filepath string) string {
	// if filepath[:7] == "http://" || filepath[:8] == "https://" {
	// 	resp, err := http.Get(filepath)
	// 	check(err)
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	check(err)
	// 	return string(body)
	// }
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

// Create creates an empty file
func Create(filepath string) {
	if !Exists(filepath) {
		Save("", filepath)
	}
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

// // ListAllDir returns list of all subdirectories in directory, scanning recursively
// func ListAllDir(dir string) []string {
// 	// what
// }

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

// Append appends the given data to the file, and creates it if it doesn't exist
func Append(filepath, data string) {
	var dat string
	if Exists(filepath) {
		dat = Open(filepath)
		dat += data
		Save(data, filepath)
	} else {
		Save(data, filepath)
	}
}

// Exists checks if the file exists
func Exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Fs is a virtual filesystem
type Fs map[string]string // filepath, contents (contents can be URL?)

// MakeFs makes a virtual filesystem
func MakeFs() (out Fs) {
	m := sync.Mutex{} // Protect against data races
	m.Lock()
	out = make(Fs)
	m.Unlock()
	return
}

// Open returns the file contents as a string
func (fs Fs) Open(filename string) string {
	m := sync.Mutex{} // Protect against data races
	m.Lock()
	if file, ok := fs[filename]; ok {
		m.Unlock()
		return file
	}
	m.Unlock()
	return ""
}

// Save saves the contents of the file to the virtual filesystem
func (fs Fs) Save(data, filepath string) {
	m := sync.Mutex{} // Protect against data races
	m.Lock()
	if data != "" && filepath != "" {
		fs[filepath] = data
	}
	m.Unlock()
}

// Exists checks if file exists
func (fs Fs) Exists(filename string) bool {
	return fs[filename] != ""
}

// Delete deletes a file
func (fs Fs) Delete(filename string) {
	m := sync.Mutex{} // Protect against data races
	m.Lock()
	delete(fs, filename)
	m.Unlock()
}

// Write is the same as Save
func (fs Fs) Write(data, filepath string) {
	fs.Save(data, filepath)
}

// List returns the list of files as a []string
func (fs Fs) List() (out []string) {
	for k := range fs {
		out = append(out, k)
	}
	return
}
