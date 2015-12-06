package main

import (
	"encoding/json"
	"testing"
)

func TestCleanupSyncDocument(t *testing.T) {
	doc := []byte(`{
    "_rev": "1-aaaaa",
    "_prop": true,
    "key": "value"
    }`)

	result, err := cleanupSyncDocument(doc)

	if err != nil {
		t.Errorf("Error executing cleanupSyncDocument: %v", err)
	}

	var jsonObject map[string]interface{}
	err = json.Unmarshal(result, &jsonObject)

	if err != nil {
		t.Errorf("Error parsing cleanupSyncDocument result", err)
	}

	if jsonObject["key"] != "value" {
		t.Fatal("assertion failed")
	}
}
