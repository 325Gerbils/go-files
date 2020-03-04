# go-files

A simple utils library for handling basic file read/write/check functions. Error handling is weak - will print an error via `log.Println(err)` and returns `""`, but will not notify the parent function that there was an error. This was deliberate - adding `if err != nil {}` checks after each file operation was annoying in the ioutil library, and I didn't want to recreate that problem. I'm  currently searching for a better (but just as simple!) way to handle errors. Perhaps a global channel?

### Installation:

```
go get -u github.com/325gerbils/go-files
```

### Usage:

```go

import (
    "fmt"
    "net/http"
    "os"

    "github.com/325gerbils/go-files"
)

// Load file contents into memory as a string:

f := files.Open("path/to/file.txt")
fmt.Println(f)

// Load() is EXACTLY the same as Open(). It's there in case I want different semantics.
f = files.Load("path/to/file.txt")
fmt.Println(f)

// Write string to file in existing folder:

data := "test write"
files.Write(data, "test.txt")

// Save() is EXACTLY the same as Write(). Again, in case I want different semantics.
files.Save(data, "test.txt")

// Write string to file in a nonexistant folder:
os.Mkdir("newfolder")
files.Write(data, "newfolder/test.txt")

// Create an empty file:
// Will not overwrite an existing file of the same name
files.Create("file/path.txt")

// List all files in a directory as []string:
fls := files.List("path/to/directory")

for _, f := range files {
    fmt.Println("Found file:", f)
}

// List all files in directory AND SUBDIRECTORIES recursively:
fls := files.ListAll("path/to/directories")

for _, f := range files {
    fmt.Println("Found file:", f)
}

// List all subdirectories in a directory (Non-recursive, see todo)
directories := files.ListDir("/path/to/directory")

for _, dir := range directories {
    fmt.Println("Found directory:", dir)
}

// Delete a file or directory:
// This is literally just a wrapper for os.RemoveAll() that prints an error. I fid myself using plain os.RemoveAll more.
files.Delete("file/path")

// Append data to a file:
data := "...the end"
files.Append("path/to/file", data)

// Check if a file exists:
if files.Exists("path/to/file") {
    fmt.Println("Exists")
} else {
    fmt.Println("Does not exist")
}

// Get uploaded file (SINGULAR!!!) from a POST request containing a multipart form:
http.HandleFunc("/processFileUpload", func(w http.ResponseWriter, r *http.Request) {
    content, name := files.GetFormFile("fileUploadFieldID", r)
    fmt.Println("File:", name, "Content:", content)
    
    files.Save(content, name)
})

// Get uploaded files from a POST request containing a multipart form:
http.HandleFunc("/processFilesUpload", func(w http.ResponseWriter, r *http.Request) {
    fls := files.GetFormFiles("fileUploadFieldID", r)
    
    // fls is type map[string]string{filename: contents} 
    for fname, fcontents := range fls {
        fmt.Println("File:", fname, "Content:", fcontents)
    
        files.Save(fcontents, fname)
    }
})

// Find all image files in a folder
// Returns files with one of the following extensions:
// "jpeg", "jpg", "png", "gif", "bmp", "ico", "tiff", "raw", "ai"

images := files.FindImages("path/to/folder")

for _, f := range images {
    fmt.Println("Found image file:", f)
}

// Find all CAD files in a folder
// Returns files with one of the following extensions:
// "dwg", "stl", "pdf", "svg", "dxf", "drw", "prt", "asm", "igs", "iges", "step", "ipt", "iam", "sldprt", "sldasm", "obj"

cad := files.FindCAD("path/to/folder")

for _, f := range cad {
    fmt.Println("Found CAD file:", f)
}

// Find all video files in a folder
// Returns files with one of the following extensions:
// "avi", "mov", "mp4", "webm", "flv", "mkv", "wmv", "m4v"

videos := files.FindVideo("path/to/folder")

for _, f := range videos {
    fmt.Println("Found video file:", f)
}

// Find all audio files in a folder
// Returns files with one of the following extensions:
// "pcm", "wav", "mp3", "aif", "m4a", "aiff", "aac", "ogg", "wma", "flac", "alac"

audio := files.FindAudio("path/to/folder")

for _, f := range audio {
    fmt.Println("Found audio file:", f)
}

// Find code files in a folder
// Returns files with one of the following extensions:
// "c", "cc", "cpp", "cs", "ino", "py", "java", "php", "h", "html", "css", "js", "go", "rb", "pl", "ts", "sql", "r", "kt", "rs", "bat", "sh"

code := files.FindCode("path/to/folder")

for _, f := range code {
    fmt.Println("Found code file:", f)
}

// Find all data files in a folder
// Returns files with one of the following extensions:
// "csv", "xls", "xlsx", "json", "xml", "db", "proto"

datafiles := files.FindDataFiles("path/to/folder")

for _, f := range datafiles {
    fmt.Println("Found data file:", f)
}

// Get the body of an HTTP request as a string:

http.HandleFunc("/getBody", func(w http.ResponseWriter, r *http.Request){
    body := files.GetBody(r)
    fmt.Println(body)
})
```

### Todo:

* Better error handling method
* Make some functions thread-safe (probably via mutex)
* Add `files.ListAllDir(path)` that scans subfolders recursively
* Change order of arguments for `files.Save()` and `files.Write()` from `data, filepath` to `filepath, data`
