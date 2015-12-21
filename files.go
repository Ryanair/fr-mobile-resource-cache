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
			localResource, err = newLocalResource(path)
			// if filepath.Ext(f.Name()) == ".json" {
			// 	localResource.FileName = path
			// 	localResource.Type = JSONType
			// } else {
			// 	localResource.FileName = path
			// 	localResource.Type = AttachmentType
			// }
		}

		fileList = append(fileList, localResource)

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

		//skip hidden files and directories
		if f.IsDir() && f.Name()[0:1] != "." {
			dirList = append(dirList, path)
		}

		return err
	})

	return dirList, err
}

func newLocalResource(path string) (LocalResource, error) {
	var localResource LocalResource

	f, err := os.Stat(path)

	if filepath.Ext(f.Name()) == ".json" {
		localResource.FileName = path
		localResource.Type = JSONType
	} else {
		localResource.FileName = path
		localResource.Type = AttachmentType
	}

	return localResource, err
}

//readFileContents returns
//file contents []byte
//documentID string
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
		if file == (LocalResource{}) || file.FileName[0:1] == "." {
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

		// TODO: check if we have a json file with the same name as the attachment,
		//and use it as a parent for the attachment, otherwise use the default template
		err = postDocument([]byte(DefaultAttachmentDoc), getDocumentID(file.FileName))
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
