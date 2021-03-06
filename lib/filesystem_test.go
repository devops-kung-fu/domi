package lib

import (
	"archive/zip"
	"os"
	"testing"
)

func TestOpenOS(t *testing.T) {
	fs := OSFS{}
	_, err := fs.Open("../__testdata__/test.txt")
	if err != nil {
		t.Error()
	}
}

func TestOpen(t *testing.T) {
	fs := mockFS{}
	_, err := fs.Open("fake.file")
	if err != nil {
		t.Error()
	}
}

func TestCopyOS(t *testing.T) {
	fs := OSFS{}
	dstFile, _ := fs.Create("../__testdata__/test1.txt")
	srcFile, _ := fs.Open("../__testdata__/test.txt")
	_, errCopy := fs.Copy(dstFile, srcFile)
	if errCopy != nil {
		t.Error()
	}
	errRemove := fs.Remove("../__testdata__/test1.txt")
	if errRemove != nil {
		t.Error()
	}
}

func TestCopy(t *testing.T) {
	fs := mockFS{}
	dstFile := fs.NewFile(0, "dstFile")
	srcFile := fs.NewFile(0, "srcFile")
	_, err := fs.Copy(dstFile, srcFile)
	if err != nil {
		t.Error()
	}
}

func TestCreate(t *testing.T) {
	fs := mockFS{}
	_, err := fs.Create("fake.file")
	if err != nil {
		t.Error()
	}
}

func TestRemove(t *testing.T) {
	fs := mockFS{}
	err := fs.Remove("fake.file")
	if err != nil {
		t.Error()
	}
}

func TestRemoveAll(t *testing.T) {
	fs := mockFS{}
	err := fs.RemoveAll("fakeDir")
	if err != nil {
		t.Error()
	}
}

func TestStatOS(t *testing.T) {
	fs := OSFS{}
	_, err := fs.Stat("../__testdata__/test.txt")
	if err != nil {
		t.Error()
	}
}

func TestStat(t *testing.T) {
	fs := mockFS{}
	_, err := fs.Stat("fake.file")
	if err != nil {
		t.Error()
	}
}

func TestWalkOS(t *testing.T) {
	fs := OSFS{}
	err := fs.Walk("../__testdata__/", func(filePath string, info os.FileInfo, err error) error {
		return nil
	})
	if err != nil {
		t.Error()
	}
}

func TestWalk(t *testing.T) {
	fs := mockFS{}
	err := fs.Walk("fakeDir/", func(filePath string, info os.FileInfo, err error) error {
		return nil
	})
	if err != nil {
		t.Error()
	}
}

func TestReadFileOS(t *testing.T) {
	fs := OSFS{}
	_, err := fs.ReadFile("../__testdata__/test.txt")
	if err != nil {
		t.Error()
	}
}

func TestReadFile(t *testing.T) {
	fs := mockFS{}
	_, err := fs.ReadFile("fake.file")
	if err != nil {
		t.Error()
	}
}

func TestWriteFileOS(t *testing.T) {
	fs := OSFS{}
	errWrite := fs.WriteFile("../__testdata__/test1.txt", []byte("Test"), 0644)
	if errWrite != nil {
		t.Error()
	}
	errRemove := fs.Remove("../__testdata__/test1.txt")
	if errRemove != nil {
		t.Error()
	}
}

func TestWriteFile(t *testing.T) {
	fs := mockFS{}
	err := fs.WriteFile("fake.file", []byte("Test"), 0644)
	if err != nil {
		t.Error()
	}
}

func TestNewFileOS(t *testing.T) {
	fs := OSFS{}
	result := fs.NewFile(13337, "../__testdata__/test1.txt")
	if result == nil {
		t.Error()
	}
}

func TestNewFile(t *testing.T) {
	fs := mockFS{}
	result := fs.NewFile(0, "fake.file")
	if result != nil {
		t.Error()
	}
}

func TestSanitizeExtractPathPass(t *testing.T) {
	err := sanitizeExtractPath("fake.file", "fakePath")
	if err != nil {
		t.Error()
	}
}

func TestSanitizeExtractPathFail(t *testing.T) {
	err := sanitizeExtractPath("../../fake.file", "fakePath")
	if err == nil {
		t.Error()
	}
}

func TestUnZip(t *testing.T) {
	fs := OSFS{}
	err := UnZip("../__testdata__/test.zip", "../__testdata__/test-zip")
	if err != nil {
		t.Error()
	}
	fs.RemoveAll("../__testdata__/test-zip")
}

func TestFindFiles(t *testing.T) {
	fs := OSFS{}
	files, err := FindFiles(fs, "../__testdata__", "(.*)\\.txt")
	if len(files) != 1 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestFileClose(t *testing.T) {
	archive, err := zip.OpenReader("../__testdata__/test.zip")
	if err != nil {
		t.Error()
	}
	defer CheckClose(&archive)
}
