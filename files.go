package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//scanResourcesDir scans the root directory from config.json
//and returns all json files, ignoring .git
func scanResourcesDir() ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(config.ResourcesDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//ignore git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		//we are interested only in json files
		//no hidden files either
		if filepath.Ext(f.Name()) == ".json" && f.Name()[0:1] != "." {
			fileList = append(fileList, path)
		}

		return err
	})

	return fileList, err
}

func readFileContents(file string) ([]byte, string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, "", err
	}

	basename := filepath.Base(file)

	fileName := strings.TrimSuffix(basename, filepath.Ext(basename))

	return b, fileName, nil
}
