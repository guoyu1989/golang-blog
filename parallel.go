package main

import (
	"crypto/md5"
	"errors"
  "flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
  numDigester = 10
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

var bounded = flag.Bool(
    "bounded", false,
    "Whether to use bounded number of digesters to walk through the directory")

func main() {
  var m map[string][md5.Size]
  if *bounded == false {
    m, err := MD5All(os.Args[1])
  } else {
    m, err := BoundedMD5All(os.Args[1])
  }
	if err != nil {
		fmt.Println(err)
	}
	for path, sum := range m {
		fmt.Printf("%x %s\n", sum, path)
	}
}

func sumFiles(root string, done <-chan struct{}) (
	<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup

		err := filepath.Walk(
			root,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					data, err := ioutil.ReadFile(path)
					select {
					case c <- result{path, md5.Sum(data), err}:
					case <-done:
					}
				}()
				select {
				case <-done:
					return errors.New("Walk is cancelled")
				default:
					return nil
				}
			})
		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()
	return c, errc
}

func walkFiles(root string, done <-chan struct{}) (
    <-chan string, <-chan error) {
  paths := make(chan string)
  errc := make(chan error, 1)
  go func() {
    defer close(paths)

    errc <- filepath.Walk(
        root,
        func(path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }
      if !info.Mode().IsRegular() {
        return nil
      }
      select {
      case paths <- path:
      case <-done:
        return errors.New("Walk cancelled")
      }
      return nil
    })
  }()
  return paths, errc
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
  for path := range paths {
    data, err := ioutil.ReadFile(path)
    select {
    case c <- result{path, md5.Sum(data), err}:
    case done:
      return
    }
  }
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(root, done)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func BoundedMD5All(root string) (map[string][md5.Size]byte, error) {
  done := make(chan struct{})
  defer close(done)

  var wg sync.WaitGroup
  wg.Add(numDigester)

  paths, errc := walkFiles(root, done)
  c := make(chan result)
  for i = 0; i < numDigester; i++ {
    go func() {
      digester(done, paths, c)
      wg.Done()
    }
  }
  go func() {
    wg.Wait()
    close(c)
  }
  m := make(map[string][md5.Size]byte)
  for r := range c {
    if r.err != nil {
      return nil, r.err
    }
    m[r.path] = r.sum
  }
  if err := <-errc; err != nil {
    return nil, err
  }
  return m, nil
}
