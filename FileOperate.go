package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Name string
	Path string
}

/*func main() {
	dir := `C:\Users\zhuji\Downloads`
	destPath := `C:\Users\zhuji\Desktop`

	mp4FileList := ReadMp4File(dir)
	for _, mp4File := range mp4FileList {
		moveMp4File(mp4File.Path, destPath)
	}
}*/

func ReadMp4File(path string) []FileInfo {
	var mp4FileList []FileInfo
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.EqualFold(filepath.Ext(info.Name()), ".mp4") {
			mp4FileList = append(mp4FileList, FileInfo{info.Name(), path})
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return mp4FileList
}

func moveMp4File(source string, destDir string) {
	sourceFileName := filepath.Base(source)
	destPath := filepath.Join(destDir, sourceFileName)
	err := os.Rename(source, destPath)
	if err != nil {
		fmt.Println(err)
	}
}
