package files

import (
	"fmt"
	"io/ioutil"
)

// Open returns the file contents as a string
func Open(filepath string) string {
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

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
