package main

import (
	"fmt"
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

//patchFiles generates a patch between the local json resource and the remote document
//and creates a new revision if there is a diff
//returns patch list, error
func patchFiles(files []string) ([]string, error) {
	var (
		returnError error
		patches     []string
	)

	for _, file := range files {
		localDocument, fileName, err := readFileContents(file)

		if err != nil {
			// log the error, don't stop execution if a file fails to read
			returnError = fmt.Errorf("Error reading file contents: %v", err)
			continue
		}

		syncDocument, _, err := getDocument(fileName)
		if err != nil {
			returnError = fmt.Errorf("Error reading sync document: %v", err)
			continue
		}

		syncDocument, err = cleanupSyncDocument(syncDocument)
		if err != nil {
			returnError = fmt.Errorf("Error cleaning up sync document: %v", err)
			continue
		}

		//if the document does not exist, or there is a difference in revisions, post it
		if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localDocument, syncDocument) {
			patch, err := diff(localDocument, syncDocument)
			if err != nil && len(syncDocument) > 0 {
				returnError = fmt.Errorf("Error generating patch: %v", err)
				continue
			}

			//print the changes in the document
			patches = append(patches, string(patch))

			err = postDocument(localDocument, fileName)
			if err != nil {
				returnError = fmt.Errorf("Error saving document: %v", err)
				continue
			}
		}
	}

	return patches, returnError
}
