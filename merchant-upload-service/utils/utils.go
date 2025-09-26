package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func TravalPath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(path, entry.Name())
			if err := TravalPath(fullPath); err != nil {
				return err
			}
		}
	} else {
		fmt.Println("文件:", path)
	}
	return nil
}
