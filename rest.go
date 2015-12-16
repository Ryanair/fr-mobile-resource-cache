package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/couchbaselabs/logg"
)

var globalHTTP = &http.Client{}

func readResource(url string) ([]byte, error) {
	logg.LogTo(TagLog, "Getting %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	document, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return document, err
}

// getDocument queries a document via sync gateway's REST API
// and returns the document contents and last revision
func getDocument(documentID string) ([]byte, string, error) {
	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + documentID

	result, err := readResource(syncEndpoint)

	var jsonObject map[string]interface{}
	err = json.Unmarshal(result, &jsonObject)

	if err != nil {
		return nil, "", err
	}

	rev, ok := jsonObject["_rev"].(string)

	if ok {
		return result, rev, nil
	}

	return nil, "", nil
}

func postDocument(document []byte, documentID string) error {
	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + documentID

	_, rev, err := getDocument(documentID)

	if rev != "" {
		syncEndpoint += "?rev=" + rev
	}

	request, err := http.NewRequest("PUT", syncEndpoint, bytes.NewReader(document))
	request.ContentLength = int64(len(document))

	logRequest(request)

	response, err := globalHTTP.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	logg.LogTo(TagLog, "%s", contents)

	return nil
}

// func createFileUploadRequest(documentID string, filePath string, contentType string) (*http.Request, error) {
// 	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + documentID
//
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
//
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = io.Copy(part, file)
//
// 	err = writer.Close()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	req, err := http.NewRequest("PUT", syncEndpoint, body)
// 	req.Header.Add("Content-Type", contentType)
//
// 	return req, err
// }

func postAttachment(fileContents []byte, parentDoc string, documentName string) error {
	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + parentDoc + "/" + documentName

	request, err := http.NewRequest("PUT", syncEndpoint, bytes.NewReader(fileContents))
	request.Header.Add("Content-Type", http.DetectContentType(fileContents))

	logRequest(request)

	response, err := globalHTTP.Do(request)

	defer response.Body.Close()

	logg.LogTo(TagLog, "Post status code: %v", response.StatusCode)

	return err
}
