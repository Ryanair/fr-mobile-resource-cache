package main

import (
	"bytes"
	"os"
	"path/filepath"
)

//JSONType specifies the local file to be of type json
var JSONType = "json"

//AttachmentType specifies the local file to be of type json
var AttachmentType = "attachment"

//LocalResource @FileName - path to the local json resource
//@Content - body of the json resource
//@ResourceID - id of the document
//@Attachment - path to the document attachment
type LocalResource struct {
	FileName   string
	Type       string
	Content    []byte
	ResourceID string
	Attachment string
}

//NewLocalResource creates a LocalResource object from a given file path
func NewLocalResource(path string) (LocalResource, error) {
	var localResource LocalResource

	f, err := os.Stat(path)

	if filepath.Ext(f.Name()) == ".json" {
		localResource.FileName = path
		localResource.Type = JSONType

		contents, documentID, err := readFileContents(path)
		if err != nil {
			return localResource, err
		}

		localResource.Content = contents
		localResource.ResourceID = documentID

	} else {
		// TODO: find the parent document, if there's one
		// localResource.FileName = path
		// localResource.Type = AttachmentType
	}

	return localResource, err
}

func (a LocalResource) compare(b LocalResource) bool {
	if &a == &b {
		return true
	}
	if a.FileName != b.FileName || a.Type != b.Type || a.ResourceID != b.ResourceID || a.Attachment != b.Attachment {
		return false
	}

	if !bytes.Equal(a.Content, b.Content) {
		return false
	}

	return true
}