package datahandle

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func GetCurrentPath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("can not get current file info"))
	}
	fmt.Println(file)
	return file
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()

}

func IsFile(path string) bool {
	return !IsDir(path)
}

func EnumAllFilePathInDir(dirs ...interface{}) []string {
	queue := list.New()
	retPaths := make([]string, 0)

	for _, dir := range dirs {
		queue.PushBack(dir.(string))
	}

	for queue.Len() != 0 {
		curDir := queue.Front().Value.(string)
		files, err := ioutil.ReadDir(curDir)
		if err != nil {
			log.Printf("Invalid path `%s`. will ignore...", curDir)
			continue
		} else {
			for _, v := range files {
				if v.IsDir() {
					fmt.Println("dir:\t", v.Name())
					queue.PushBack(v.Name())
				} else {
					fmt.Println("file:\t", v.Name())
					retPaths = append(retPaths, v.Name())
				}
			}
		}
	}

	return retPaths
}
