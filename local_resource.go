package main

import (
	"bytes"

	"github.com/couchbaselabs/logg"
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

//NewLocalDocument creates a LocalResource object from a given file path
func NewLocalDocument(path string, resourceList *[]LocalResource) error {
	var localResource LocalResource

	localResource.FileName = path
	localResource.Type = JSONType

	contents, documentID, err := readFileContents(path)
	if err != nil {
		return err
	}

	localResource.Content = contents
	localResource.ResourceID = documentID

	*resourceList = append(*resourceList, localResource)

	return err
}

//NewLocalAttachment ..
func NewLocalAttachment(path string, resourceList *[]LocalResource) error {
	documentID := getDocumentID(path)
	result := ResourceIndexAt(resourceList, documentID)

	if result == -1 {
		logg.LogTo(TagLog, "Parent document not found, creating a new one ...")
		var localResource LocalResource
		localResource.FileName = path
		localResource.Type = AttachmentType
		localResource.Attachment = path
		localResource.Content = []byte(DefaultAttachmentDoc)
		localResource.ResourceID = documentID

		*resourceList = append(*resourceList, localResource)
	} else {
		logg.LogTo(TagLog, "Parent document found at %d", result)
		(*resourceList)[result].Attachment = path
	}

	return nil
}

//ResourceIndexAt searches for a specific documentID in a LocalResource slice and returns the index of the found element
//-1 if no result is found
func ResourceIndexAt(resourceList *[]LocalResource, documentID string) int {
	for i, resource := range *resourceList {
		if resource.ResourceID == documentID {
			return i
		}
	}

	return -1
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
