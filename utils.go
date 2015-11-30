package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/couchbaselabs/logg"
)

func getRootNode(m map[string][]interface{}) string {
	for i := range m {
		return i
	}

	return ""
}

func logRequest(request *http.Request) {
	if DEBUG {
		log, _ := httputil.DumpRequest(request, true)
		logg.Log("%s", log)
	}
}

func cleanupSyncDocument(syncDocument []byte) ([]byte, error) {
	if len(syncDocument) == 0 {
		return nil, nil
	}

	document := make(map[string]interface{})
	err := json.Unmarshal(syncDocument, &document)

	for i := range document {
		if m, err := regexp.MatchString("_([a-z]+)", i); m == true && err == nil {
			delete(document, i)
		}
	}

	result, err := json.Marshal(document)

	return result, err
}
