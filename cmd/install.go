package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	tmpZipFile = "tmp.zip"
	dbFile     = "dingo.db"
)

func CheckInstall() bool {
	_, err := os.Stat(dbFile)
	return err == nil
}

func ExtractBundleBytes() error {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(siteContent))
	b, _ := ioutil.ReadAll(decoder)
	ioutil.WriteFile(tmpZipFile, b, os.ModePerm)
	reader, err := zip.OpenReader(tmpZipFile)
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join("", file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func Install() {
	println("Installation complete")
}
