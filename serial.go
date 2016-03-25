package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	for path, sum := range m {
		fmt.Printf("%x %s\n", sum, path)
	}
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			m[path] = md5.Sum(data)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return m, nil
}
