package utils

import (
	"os"
	"time"
	"io/ioutil"
	"crypto/md5"
	"fmt"
)

type DynaFile struct {
	Content []byte
	Md5sum  string
}

var (
	contents map[string]*DynaFile
)

func readFiles(files []string) {
	mtimes := make([]time.Time, len(files))

	for {
		for i, file := range files {
			fi, err := os.Lstat(file)
			if err != nil {
				continue
			}
			mtime := fi.ModTime()
			if mtimes[i].Before(mtime) {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					continue
				}
				h := md5.New()
				contents[file] = &DynaFile{b, fmt.Sprintf("%X", h.Sum(b))}
				mtimes[i] = mtime
			}
		}
		time.Sleep(2*time.Minute)
	}
}

func StartFilesThreads(files []string) {
	contents = make(map[string]*DynaFile, len(files))
	for _, file := range files {
		contents[file] = nil
	}
	go readFiles(files)
}

func GetFileContent(filename string) ([]byte, bool) {
	if f, ok := contents[filename]; !ok {
		return nil, false
	} else {
		return f.Content, ok
	}
}

func GetFile(filename string) (*DynaFile) {
	if f, ok := contents[filename]; !ok {
		return nil
	} else {
		return f
	}
}
