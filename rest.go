package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/couchbaselabs/logg"
)

var globalHTTP = &http.Client{}

func readResource(url string) []byte {
	if DEBUG {
		logg.Log("Getting %s\n", url)
	}

	res, err := http.Get(url)
	if err != nil {
		logg.LogPanic("Error parsing %s: %v", url, err.Error())
	}

	document, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		logg.LogPanic("Error reading data: %v", err.Error())
	}

	return document
}

// getDocument queries a document via sync gateway's REST API
// and returns the document contents and last revision
// panics if the document does not exist
// todo: don't panic, return nil
func getDocument(documentID string) ([]byte, string) {
	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + documentID

	result := readResource(syncEndpoint)

	var jsonObject map[string]interface{}
	err := json.Unmarshal(result, &jsonObject)

	if err != nil {
		logg.LogPanic("Error parsing document: %v", err)
		return nil, ""
	}

	rev, ok := jsonObject["_rev"].(string)

	if ok {
		return result, rev
	}

	return nil, ""
}

func postDocument(document []byte, documentID string) {
	var syncEndpoint = config.SyncURL + "/" + config.Bucket + "/" + documentID

	_, rev := getDocument(documentID)

	if rev != "" {
		syncEndpoint += "?rev=" + rev
	}

	request, err := http.NewRequest("PUT", syncEndpoint, bytes.NewReader(document))
	request.ContentLength = int64(len(document))

	logRequest(request)

	response, err := globalHTTP.Do(request)

	if err != nil {
		logg.LogPanic("Error saving document: %v", err)
		return
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		logg.LogPanic("%v", err)
		return
	}

	logg.Log("%s", contents)

}
