package main

import (
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

		if f.IsDir() {
			dirList = append(dirList, path)
		}

		return err
	})

	return dirList, err
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
		patch, err := file.updateDoc()
		if err != nil {
			returnError = err
			continue
		}
		patches = append(patches, patch)
	}

	return patches, returnError
}

func (localResource LocalResource) updateDoc() (string, error) {
	var patch []byte
	syncDocument, rev, err := getDocument(localResource.ResourceID)
	if err != nil {
		return "", err
	}

	//if the document does not exist, or there is a difference in revisions, post it
	if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localResource.Content, syncDocument) {
		patch, err = diff(localResource.Content, syncDocument)

		newRev, err := postDocument(localResource.Content, localResource.ResourceID)
		if err != nil {
			return "", err
		}

		//the document has an attachment, post it again, as it looses reference to the file after every document update
		if localResource.Attachment != "" {
			err = updateAttachment(localResource.Attachment, newRev)
			if err != nil {
				return "", err
			}
		}

		logg.LogTo(TagLog, "Patch: %s; rev: %s, newRev: %s", string(patch), rev, newRev)
	} else if localResource.Attachment != "" {
		//the document hasn't updated but we need to check the attachment if there's one
		needsUpdate, err := attachmentNeedsUpdate(localResource.Attachment)
		if needsUpdate && err == nil {
			err = updateAttachment(localResource.Attachment, rev)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return string(patch), err
}

func updateAttachment(path, rev string) error {
	body, documentID, err := readFileContents(path)
	err = postAttachment(body, documentID, filepath.Base(path)+"?rev="+rev)

	return err
}

func attachmentNeedsUpdate(path string) (bool, error) {
	body, documentID, err := readFileContents(path)
	if err != nil {
		return false, err
	}

	remoteAttachmentURL := config.SyncURL + "/" + config.Bucket + "/" + documentID + "/" + filepath.Base(path)
	remoteFile, err := readResource(remoteAttachmentURL)
	if err != nil {
		return false, err
	}

	crc := crc32.ChecksumIEEE(body)
	crcRemote := crc32.ChecksumIEEE(remoteFile)

	return crc != crcRemote, nil
}
