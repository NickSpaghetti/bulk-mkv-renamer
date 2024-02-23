package utils

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type BulkFile struct {
	Name      string
	Path      string
	Ext       string
	ParentDir string
}

func GetFiles(path string) ([]BulkFile, error) {
	var bulkFiles []BulkFile
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		bulkFiles = append(bulkFiles, BulkFile{
			Name:      info.Name(),
			Path:      path,
			Ext:       info.Name(),
			ParentDir: filepath.Dir(path),
		})
		fmt.Println(info.Name())
		return nil
	})
	return bulkFiles, err
}
