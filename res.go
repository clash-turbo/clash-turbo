package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func delFile(path, fileName string) {
	err := os.Remove(filepath.Join(path, fileName))
	if err != nil {
		return
	}
}

func workPathFileExists(fileName string) bool {
	filePath := filepath.Join(getWorkPath(), fileName)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
func exePathFileExists(fileName string) bool {
	filePath := filepath.Join(getExePath(), fileName)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func makeTextFile(content, path, fileName string) {
	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		fmt.Println("Error creating file:", fileName, "--------", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("文件创建失败")
		return
	}
}

func getExePath() string {
	currentDir, _ := os.Getwd()

	if _exePath == "" {
		return currentDir
	} else {
		return _exePath
	}

}

func getWorkPath() string {
	return filepath.Join(getExePath(), _workPath)
}
