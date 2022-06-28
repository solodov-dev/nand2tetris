package main

import (
	"log"
	"strings"
)

func PathInfo(path string) (name, dir string, isFile bool) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	splitted := strings.Split(path, "/")
	lastIndex := len(splitted) - 1
	nameExt := strings.Split(splitted[lastIndex], ".")
	name = nameExt[0]
	isDir := len(nameExt) == 1
	if isDir {
		dir = path + "/"
	} else {
		dir = strings.Join(splitted[:lastIndex], "/") + "/"
	}
	isFile = !isDir
	return
}

func HandleError(err error, msg string) {
  if err != nil {
    log.Fatalf(msg)
  }
}
