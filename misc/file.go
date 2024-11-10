package misc

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var usedFiles = make(map[string]*os.File)

var fileMutexes = make(map[string]*sync.Mutex)
var fileMutexesLock sync.Mutex

func OpenFile(directory string, fileName string) (*os.File, *sync.Mutex) {
	filePath := filepath.Join(directory, fileName)
	fileMutex := getFileMutex(filePath)

	if file, ok := usedFiles[filePath]; ok {
		return file, fileMutex
	}

	//Create directory
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	//Create / Open file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
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
func Anew(endpoints []string, fileName string, directory string, add bool) []string {

	file, m := OpenFile(directory, fileName)
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

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Check for duplicates and accumulate unique lines
	for _, line := range endpoints {

		if lines[line] {
			continue
		}
		// add the line to the map so we don't get any duplicates from stdin
		lines[line] = true
		//Append new lines to the file
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

	f, m := OpenFile(directory, fileName)
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
				log.Fatal(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err := os.WriteFile(filePath, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

}

// Read file
func Read(fileName string) ([]string, error) {
	file, m := OpenFile("./", fileName)
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

// Append file
func Write(text string, fileName string, directory string) {
	f := strings.Replace(fileName, "/", "_", -1)

	file, m := OpenFile(directory, f)
	m.Lock()
	defer m.Unlock()

	file.WriteString(text + "\n")
}

// Append slice to a file
func WriteAll(text []string, fileName string, directory string) {
	if len(text) == 0 {
		return
	}

	f := strings.Replace(fileName, "/", "_", -1)

	file, m := OpenFile(directory, f)
	m.Lock()
	defer m.Unlock()

	file.WriteString(strings.Join(text, "\n") + "\n")
}

func RemoveFile(file string) {
	os.Remove(file)
}

func getFileMutex(filePath string) *sync.Mutex {
	fileMutexesLock.Lock()
	defer fileMutexesLock.Unlock()

	if _, ok := fileMutexes[filePath]; !ok {
		fileMutexes[filePath] = &sync.Mutex{}
	}
	return fileMutexes[filePath]
}
