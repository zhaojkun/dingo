package model

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type File struct {
	os.FileInfo
	Url     string
	ModTime *time.Time
}

// Checks if directory is a child directory of base, make sure that GetFileList won't
// read any folder other than the upload folder.
func CheckSafe(directory string, base string) bool {
	directory = path.Clean(directory)
	dirs := strings.Split(directory, "/")
	return dirs[0] == base
}

func GetFileList(directory string) []*File {
	files := make([]*File, 0)
	fileInfoList, _ := ioutil.ReadDir(directory)
	for i := len(fileInfoList) - 1; i >= 0; i-- {
		if fileInfoList[i].Name() == ".DS_Store" {
			continue
		}
		file := new(File)
		file.FileInfo = fileInfoList[i]
		file.Url = path.Join(directory, fileInfoList[i].Name())
		t := fileInfoList[i].ModTime()
		file.ModTime = &t
		files = append(files, file)
	}
	return files
}

func RemoveFile(path string) error {
	return os.RemoveAll(path)
}

func CreateFilePath(dir string, name string) string {
	os.MkdirAll(dir, os.ModePerm)
	return path.Join(dir, name)
}
