package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//LocalResource represents the local document
//can be json or attachment
type LocalResource struct {
	FileName string
	Type     string
}

//JSONType specifies the local file to be of type json
var JSONType = "json"

//AttachmentType specifies the local file to be of type json
var AttachmentType = "attachment"

//scanResourcesDir scans the root directory from config.json
//and returns all json files, ignoring .git
func scanResourcesDir() ([]LocalResource, error) {
	fileList := []LocalResource{}
	err := filepath.Walk(config.ResourcesDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//ignore git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		var localResource LocalResource

		//skip hidden files and directories
		if !f.IsDir() && f.Name()[0:1] != "." {
			if filepath.Ext(f.Name()) == ".json" {
				localResource.FileName = path
				localResource.Type = JSONType
			} else {
				localResource.FileName = path
				localResource.Type = AttachmentType
			}
		}

		fileList = append(fileList, localResource)

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
func patchFiles(files []LocalResource) ([]string, error) {
	var (
		returnError error
		patches     []string
	)

	for _, file := range files {
		if file == (LocalResource{}) {
			continue
		}

		if file.Type == JSONType {
			patch, err := updateJSONDoc(file)
			if err != nil {
				returnError = err
				continue
			}

			patches = append(patches, patch)
		} else {
			err := updateAttachment(file)
			if err != nil {
				returnError = err
			}
		}
	}

	return patches, returnError
}

func updateJSONDoc(file LocalResource) (string, error) {
	var (
		returnError error
		patch       []byte
	)

	localDocument, fileName, err := readFileContents(file.FileName)

	if err != nil {
		// log the error, don't stop execution if a file fails to read
		returnError = fmt.Errorf("Error reading file contents: %v", err)
	}

	syncDocument, _, err := getDocument(fileName)
	if err != nil {
		returnError = fmt.Errorf("Error reading sync document: %v", err)
	}

	syncDocument, err = cleanupSyncDocument(syncDocument)
	if err != nil {
		returnError = fmt.Errorf("Error cleaning up sync document: %v", err)
	}

	//if the document does not exist, or there is a difference in revisions, post it
	if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localDocument, syncDocument) {
		patch, err = diff(localDocument, syncDocument)
		if err != nil && len(syncDocument) > 0 {
			returnError = fmt.Errorf("Error generating patch: %v", err)
		}

		err = postDocument(localDocument, fileName)
		if err != nil {
			returnError = fmt.Errorf("Error saving document: %v", err)
		}
	}

	return string(patch), returnError
}

func updateAttachment(file LocalResource) error {
	var returnError error

	_, fileName, err := readFileContents(file.FileName) // TODO: handle the local file

	if err != nil {
		// log the error, don't stop execution if a file fails to read
		returnError = fmt.Errorf("Error reading file contents: %v", err)
	}

	syncDocument, _, err := getDocument(fileName)
	if err != nil {
		returnError = fmt.Errorf("Error reading sync document: %v", err)
	}

	dummyDoc := []byte(`{
			"channels": "ch_v1"
		}`)

	if len(syncDocument) == 0 || syncDocument == nil {
		err = postDocument(dummyDoc, fileName)
		if err != nil {
			returnError = fmt.Errorf("Error saving document: %v", err)
		}
	}

	return returnError
}
