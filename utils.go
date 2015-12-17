package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	document := make(map[string]interface{})
	err := json.Unmarshal(syncDocument, &document)

	if err != nil {
		return nil, err
	}

	for i := range document {
		if m, err := regexp.MatchString("_([a-z]+)", i); m == true && err == nil {
			delete(document, i)
		}
	}

	return json.Marshal(document)
}

func getReadSyncEndpoint() string {
	return config.SyncURL + "/" + config.Bucket + "/"
}

func getWriteSyncEndpoint() string {
	if config.Username != "" && config.Password != "" {
		rawurl := config.SyncURL + "/" + config.Bucket + "/"

		url, err := url.Parse(rawurl)

		if err != nil {
			logg.LogFatal("%s", err)
		}

		writeURL := url.Scheme + "://" + config.Username + ":" + config.Password + "@" + url.Host + url.Path

		return writeURL
	}

	return config.SyncURL + "/" + config.Bucket + "/"

}
