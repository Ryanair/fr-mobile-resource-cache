package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/couchbaselabs/logg"
)

//scanResourcesDir scans the root directory from config.json
//and returns all json and binary files, ignoring .git
func scanResourcesDir() ([]string, error) {
	var fileList []string
	err := filepath.Walk(config.ResourcesDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//ignore git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		//skip hidden files and directories
		if !f.IsDir() && f.Name()[0:1] != "." {
			fileList = append(fileList, path)
		}

		return err
	})

	return fileList, err
}

func getLocalResources() ([]LocalResource, error) {
	files, err := scanResourcesDir()

	var result []LocalResource

	//first pass, get all json documents
	for _, file := range files {
		if filepath.Ext(file) == ".json" {
			err := NewLocalDocument(file, &result)
			if err != nil {
				continue
			}
		}
	}

	//second pass, get all attachments
	for _, file := range files {
		if filepath.Ext(file) != ".json" {
			err := NewLocalAttachment(file, &result)
			if err != nil {
				continue
			}
		}
	}

	return result, err
}

func findFile(fileID string) (string, error) {
	var filename string
	err := filepath.Walk(config.ResourcesDir, func(path string, f os.FileInfo, err error) error {
		//ignore git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		if getDocumentID(path) == fileID {
			filename = path
		}

		return err
	})

	return filename, err
}

func getDirectories() ([]string, error) {
	var dirList []string
	err := filepath.Walk(config.ResourcesDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//ignore git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		//skip hidden files and directories
		if f.IsDir() && f.Name()[0:1] != "." {
			dirList = append(dirList, path)
		}

		return err
	})

	return dirList, err
}

//readFileContents @file - path to input file
//[]byte - file contents
//string - documentID
//error
func readFileContents(file string) ([]byte, string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, "", err
	}

	fileName := getDocumentID(file)

	return b, fileName, nil
}

func getDocumentID(fileName string) string {
	basename := filepath.Base(fileName)

	return strings.TrimSuffix(basename, filepath.Ext(basename))
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
		if file.compare(LocalResource{}) || file.FileName[0:1] == "." {
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
			logg.LogTo(TagLog, "Trying to update %v", file)
			err := updateAttachment(file)
			if err != nil {
				returnError = err
			}
		}
	}

	return patches, returnError
}

func updateJSONDoc(localResource LocalResource) (string, error) {
	var (
		returnError error
		patch       []byte
	)

	syncDocument, _, err := getDocument(localResource.ResourceID)
	if err != nil {
		returnError = fmt.Errorf("Error reading sync document: %v", err)
	}

	syncDocument, err = cleanupSyncDocument(syncDocument)
	if err != nil {
		returnError = fmt.Errorf("Error cleaning up sync document: %v", err)
	}

	//if the document does not exist, or there is a difference in revisions, post it
	if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localResource.Content, syncDocument) {
		patch, err = diff(localResource.Content, syncDocument)
		if err != nil && len(syncDocument) > 0 {
			returnError = fmt.Errorf("Error generating patch: %v", err)
		}

		err = postDocument(localResource.Content, localResource.ResourceID)
		if err != nil {
			returnError = fmt.Errorf("Error saving document: %v", err)
		}
	}

	return string(patch), returnError
}

func updateAttachment(file LocalResource) error {
	var returnError error

	syncDocument, lastRev, err := getDocument(getDocumentID(file.FileName))

	if err != nil {
		returnError = fmt.Errorf("Error reading sync document: %v", err)
	}

	fileBody, docID, err := readFileContents(file.FileName)

	//create a new document
	if len(syncDocument) == 0 || syncDocument == nil {
		if err != nil {
			return fmt.Errorf("Error reading attachment: %s", err)
		}

		var postDocContents []byte
		attachmentDoc, err := findFile(getDocumentID(file.FileName))
		if attachmentDoc != "" {
			postDocContents, _, err = readFileContents(attachmentDoc)
		} else {
			postDocContents = []byte(DefaultAttachmentDoc)
		}

		err = postDocument(postDocContents, getDocumentID(file.FileName))
		if err != nil {
			return fmt.Errorf("Error saving document: %v", err)
		}

		// TODO: do that only for new files, otherwise extract from first document read
		_, rev, err := getDocument(getDocumentID(file.FileName))
		if err != nil {
			returnError = err
		}

		if rev != "" {
			err = postAttachment(fileBody, getDocumentID(file.FileName), filepath.Base(file.FileName)+"?rev="+rev)
		}
	} else {
		// generate file checksum
		remoteAttachmentURL := config.SyncURL + "/" + config.Bucket + "/" + docID + "/" + filepath.Base(file.FileName)
		remoteFile, err := readResource(remoteAttachmentURL)
		if err != nil {
			returnError = err
		}

		crc := crc32.ChecksumIEEE(fileBody)
		crcRemote := crc32.ChecksumIEEE(remoteFile)

		// update an exsisting document
		if crc != crcRemote && lastRev != "" {
			err = postAttachment(fileBody, getDocumentID(file.FileName), filepath.Base(file.FileName)+"?rev="+lastRev)
		}
	}

	return returnError
}
