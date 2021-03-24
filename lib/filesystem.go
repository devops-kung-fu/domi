package lib

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type fileSystem interface {
	Open(name string) (file, error)
	Copy(dst io.Writer, src io.Reader) (int64, error)
	Create(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type file interface {
	io.Closer
	io.Writer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// OSFS implements fileSystem using the local disk.
type OSFS struct{}

// Open - Open File
func (OSFS) Open(name string) (file, error) { return os.Open(name) }

// Copy - Copy File
func (OSFS) Copy(dst io.Writer, src io.Reader) (int64, error) { return io.Copy(dst, src) }

// Create - Create File
func (OSFS) Create(name string) (file, error) { return os.Create(name) }

// Stat - Stat File
func (OSFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// Walk - Walk Path
func (OSFS) Walk(root string, walkFn filepath.WalkFunc) error { return filepath.Walk(root, walkFn) }

// ReadFile - Reads a File
func (OSFS) ReadFile(filename string) ([]byte, error) { return ioutil.ReadFile(filename) }

// WriteFile - Writes to a File
func (OSFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

type mockFS struct{}

func (mockFS) Open(name string) (file, error)                                 { return nil, nil }
func (mockFS) Copy(dst io.Writer, src io.Reader) (int64, error)               { return 100, nil }
func (mockFS) Create(name string) (file, error)                               { return nil, nil }
func (mockFS) Stat(name string) (os.FileInfo, error)                          { return nil, nil }
func (mockFS) Walk(root string, walkFn filepath.WalkFunc) error               { return nil }
func (mockFS) ReadFile(filename string) ([]byte, error)                       { return []byte(`Test String`), nil }
func (mockFS) WriteFile(filename string, data []byte, perm os.FileMode) error { return nil }

// DownloadFile - Download a file from a URL
func DownloadFile(fs fileSystem, url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}
	thisUUID := getUUID()
	fileName := fmt.Sprintf("/domi/%s.zip", thisUUID)
	file, err := fs.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = fs.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Downloaded %s as %s\n", url, fileName)
	return thisUUID, nil
}

func sanitizeExtractPath(filePath string, destination string) error {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, destination) {
		return fmt.Errorf("%s: illegal file path", filePath)
	}
	return nil
}

// UnZip - Extracts a zip archive
func UnZip(source string, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		err = sanitizeExtractPath(file.Name, destination)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	log.Printf("Extracted %s into %s", source, destination)
	return nil
}

// FindFiles - Recursively search for files matching a pattern.
func FindFiles(fs fileSystem, root string, re string) ([]string, error) {
	libRegEx, e := regexp.Compile(re)
	if e != nil {
		return nil, e
	}
	var files []string
	e = fs.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			files = append(files, filePath)
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, nil
}
