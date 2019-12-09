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
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	files, err := ioutil.ReadDir(dir)
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
	files, err := ioutil.ReadDir(dir)
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

// Delete deletes the files at the given path
func Delete(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	}
}

// Remove removes the files at the given path
func Remove(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	}
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

// GetFormFile returns the file (contents and name) associated with the form field
func GetFormFile(fname string, r *http.Request) (content, name string) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	filename := handler.Filename
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	return string(fileBytes), filename
}

// GetFormFiles returns the list of files associated with the form field
func GetFormFiles(fname string, r *http.Request) ([]string, []string) {
	r.ParseMultipartForm(32 << 20)
	fhs := r.MultipartForm.File[fname]
	files := make([]string, len(fhs))
	filenames := make([]string, len(fhs))
	for i, fh := range fhs {
		filename := fh.Filename
		f, err := fh.Open()
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
		}
		files[i] = string(fileBytes)
		filenames[i] = filename
	}
	return files, filenames
}

// FindImages returns list of image files in dir.
// Image extensions: jpeg, jpg, png, gif, bmp, ico, tiff, raw, ai
func FindImages(path string) (out []string) {
	imgexts := []string{"jpeg", "jpg", "png", "gif", "bmp", "ico", "tiff", "raw", "ai"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range imgexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// FindCAD returns list of CAD files in dir.
// CAD extensions: dwg, stl, pdf, svg, dxf, drw, prt, asm, igs, iges, step,
// ipt, iam, sldprt, sldasm, obj
func FindCAD(path string) (out []string) {
	cadexts := []string{"dwg", "stl", "pdf", "svg", "dxf", "drw", "prt", "asm", "igs", "iges", "step", "ipt", "iam", "sldprt", "sldasm", "obj"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range cadexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// FindVideo returns list of video files in dir.
// Video extensions: avi, mov, mp4, webm, flv, mkv, wmv, m4v
func FindVideo(path string) (out []string) {
	videoexts := []string{"avi", "mov", "mp4", "webm", "flv", "mkv", "wmv", "m4v"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range videoexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// FindAudio returns list of audio files in dir.
// Audio extensions: pcm, wav, mp3, aif, m4a, aiff, aac, ogg, wma, flac, alac
func FindAudio(path string) (out []string) {
	audioexts := []string{"pcm", "wav", "mp3", "aif", "m4a", "aiff", "aac", "ogg", "wma", "flac", "alac"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range audioexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// FindCode returns list of code files in dir.
// Code extensions: c, cc, cpp, cs, ino, py, java, php, h, html, css, js,
// go, rb, pl, ts, sql, r, kt, rs, bat, sh
func FindCode(path string) (out []string) {
	codeexts := []string{"c", "cc", "cpp", "cs", "ino", "py", "java", "php", "h", "html", "css", "js", "go", "rb", "pl", "ts", "sql", "r", "kt", "rs", "bat", "sh"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range codeexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// FindDataFiles returns list of data files in dir.
// Data files extensions: csv, xls, xlsx, json, xml, db, proto
func FindDataFiles(path string) (out []string) {
	dataexts := []string{"csv", "xls", "xlsx", "json", "xml", "db", "proto"}
	filesList := ListAll(path)
	for _, file := range filesList {
		for _, ext := range dataexts {
			if filepath.Ext(file) == ext {
				out = append(out, file)
			}
		}
	}
	return
}

// GetBody returns the body (as string) of an *http.Request
func GetBody(r *http.Request) string {
	n, err := ioutil.ReadAll(r.Body)
	check(err)
	n = n[1 : len(n)]
	return string(n)
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
