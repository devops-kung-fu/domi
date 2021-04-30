package lib

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileSystem - Interface
type FileSystem interface {
	Open(name string) (File, error)
	Copy(dst io.Writer, src io.Reader) (int64, error)
	Create(name string) (File, error)
	Remove(name string) error
	RemoveAll(name string) error
	Stat(name string) (os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	NewFile(fd uintptr, name string) File
}

// File interface
type File interface {
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
func (OSFS) Open(name string) (File, error) { return os.Open(name) }

// Copy - Copy File
func (OSFS) Copy(dst io.Writer, src io.Reader) (int64, error) { return io.Copy(dst, src) }

// Create - Create File
func (OSFS) Create(name string) (File, error) { return os.Create(name) }

// Remove - Remove File
func (OSFS) Remove(name string) error { return os.Remove(name) }

// RemoveAll - Remove Directory
func (OSFS) RemoveAll(name string) error { return os.RemoveAll(name) }

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

// NewFile - Creates a new File
func (OSFS) NewFile(fd uintptr, name string) File {
	return os.NewFile(fd, name)
}

type mockFS struct{}

func (mockFS) Open(name string) (File, error)                                 { return nil, nil }
func (mockFS) Copy(dst io.Writer, src io.Reader) (int64, error)               { return 100, nil }
func (mockFS) Create(name string) (File, error)                               { return os.NewFile(0, "fake"), nil }
func (mockFS) Remove(name string) error                                       { return nil }
func (mockFS) RemoveAll(name string) error                                    { return nil }
func (mockFS) Stat(name string) (os.FileInfo, error)                          { return nil, nil }
func (mockFS) Walk(root string, walkFn filepath.WalkFunc) error               { return nil }
func (mockFS) ReadFile(filename string) ([]byte, error)                       { return []byte(`Test String`), nil }
func (mockFS) WriteFile(filename string, data []byte, perm os.FileMode) error { return nil }
func (mockFS) NewFile(fd uintptr, name string) File                           { return nil }

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

	defer func() {
		err = archive.Close()
	}()

	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			err = reader.Close()
		}()
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

		defer func() {
			err = writer.Close()
		}()

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	log.Printf("Extracted %s into %s", source, destination)
	return nil
}

// FindFiles - Recursively search for files matching a pattern.
func FindFiles(fs FileSystem, root string, re string) ([]string, error) {
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
