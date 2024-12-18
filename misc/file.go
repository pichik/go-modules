package misc

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var usedFiles = make(map[string]*os.File)

var fileMutexes = make(map[string]*sync.Mutex)
var fileMutexesLock sync.Mutex

func openFile(directory string, fileName string, writeType int) (*os.File, *sync.Mutex) {

	//Check directory and file name length to prevent errors
	if len(directory) > 250 {
		directory = fmt.Sprintf("directory-too-long/%s", generateRandomString(10))
	}
	if len(fileName) > 250 {
		fileName = fmt.Sprintf("filename-too-long-%s", generateRandomString(10))
	}

	filePath := filepath.Join(directory, fileName)
	fileMutex := getFileMutex(filePath)

	if file, ok := usedFiles[filePath]; ok {
		return file, fileMutex
	}

	//Create directory
	if directory != "" {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			PrintError("Create dir", err)
		}
	}

	//Create / Open file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|writeType, 0644)
	if err != nil {
		PrintError("Open file", err)
	}

	usedFiles[filePath] = file
	return file, fileMutex
}

func CloseAllFiles() {
	for _, file := range usedFiles {
		file.Close()
	}
}

// Used tomnomnom anew
// Compare strings with file lines and return them if they are not already in.
// True for writing them to file
// Variadic parameter 'endpoints ...string' accept single string or slice (have to be used as last parameter)
// With passing slice is required to use '...' Anew("file", "dir", true, newUrls...)
func Anew(fileName string, directory string, add bool, endpoints ...string) []string {
	fileName = strings.Replace(fileName, "/", "_", -1)

	file, m := openFile(directory, fileName, os.O_APPEND)
	m.Lock()
	defer m.Unlock()

	file.Seek(0, io.SeekStart)
	sc := bufio.NewScanner(file)

	lines := make(map[string]bool)
	unique := []string{}
	// Read existing lines from the file
	for sc.Scan() {
		lines[sc.Text()] = true
	}

	// Create the writer once before the loop
	var writer *bufio.Writer
	if add {
		writer = bufio.NewWriter(file)
		defer writer.Flush()
	}

	// Check for duplicates and accumulate unique lines
	for _, line := range endpoints {
		if lines[line] {
			continue
		}
		// Add the line to the map to avoid duplicates
		lines[line] = true
		// Append new lines to the file
		if add {
			writer.WriteString(line + "\n")
		}
		unique = append(unique, line)
	}

	return unique
}

// Remove matched line from file
func RemoveLine(line string, fileName string, directory string) {
	filePath := filepath.Join(directory, fileName)

	f, m := openFile(directory, fileName, os.O_TRUNC)
	m.Lock()
	defer m.Unlock()

	//Get to the beginning of the file
	f.Seek(0, io.SeekStart)

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != line {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				PrintError("File scanner", err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				PrintError("File scanner", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		PrintError("File scanner", err)
	}

	err := os.WriteFile(filePath, buf.Bytes(), 0666)
	if err != nil {
		PrintError("File scanner", err)
	}

}

// Read file
func Read(fileName string) ([]string, error) {
	file, m := openFile("./", fileName, 0)
	m.Lock()
	defer m.Unlock()

	lines := []string{}

	//Get to the beginning of the file
	file.Seek(0, io.SeekStart)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, nil
}

// Append slice to a file
func Append(fileName string, directory string, text ...string) {
	if len(text) == 0 {
		return
	}

	f := strings.Replace(fileName, "/", "_", -1)

	file, m := openFile(directory, f, os.O_APPEND)
	m.Lock()
	defer m.Unlock()

	for _, line := range text {
		file.WriteString(line + "\n")
	}
}

// Ovewrite entire file
func Overwrite(fileName string, directory string, text ...string) {
	if len(text) == 0 {
		return
	}

	f := strings.Replace(fileName, "/", "_", -1)

	file, m := openFile(directory, f, os.O_TRUNC)
	m.Lock()
	defer m.Unlock()

	for _, line := range text {
		file.WriteString(line + "\n")
	}
}

func RemoveFile(file string) {
	os.Remove(file)
}

// Each file have its own mutex
func getFileMutex(filePath string) *sync.Mutex {
	fileMutexesLock.Lock()
	defer fileMutexesLock.Unlock()

	if _, ok := fileMutexes[filePath]; !ok {
		fileMutexes[filePath] = &sync.Mutex{}
	}
	return fileMutexes[filePath]
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomByte, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randomByte.Int64()]
	}
	return string(b)
}
