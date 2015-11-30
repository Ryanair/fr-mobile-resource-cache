package main

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	"net/http"
	"net/http/httputil"
	"regexp"
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

func cleanupSyncDocument(syncDocument []byte) []byte {
	if len(syncDocument) == 0 {
		return nil
	}

	document := make(map[string]interface{})
	err := json.Unmarshal(syncDocument, &document)

	if err != nil {
		logg.LogPanic("%v", err)
	}

	for i := range document {
		if m, err := regexp.MatchString("_([a-z]+)", i); m == true && err == nil {
			delete(document, i)
		}
	}

	result, err := json.Marshal(document)

	if err != nil {
		logg.LogPanic("%v", err)
	}

	return result
}
