package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/couchbaselabs/logg"
)

func logRequest(request *http.Request) {
	log, _ := httputil.DumpRequest(request, true)
	logg.LogTo(TagLog, "%s", log)
}

func cleanupSyncDocument(syncDocument []byte) ([]byte, error) {
	if len(syncDocument) == 0 {
		return nil, nil
	}

	cleanDocument := make(map[string]interface{})
	err := json.Unmarshal(syncDocument, &cleanDocument)

	if err != nil {
		return nil, err
	}

	for i := range cleanDocument {
		if m, err := regexp.MatchString("^_([a-z]+)", i); m == true && err == nil {
			delete(cleanDocument, i)
		}
	}

	return json.Marshal(cleanDocument)
}

func getSyncEndpoint() string {
	return config.SyncURL + "/" + config.Bucket + "/"
}
