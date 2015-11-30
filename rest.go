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
	if DEBUG {
		logg.LogTo(TagLog, "Getting %s\n", url)
	}

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
// panics if the document does not exist
// todo: don't panic, return nil
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

	if DEBUG {
		logg.LogTo(TagLog, "%s", contents)
	}

	return nil
}
