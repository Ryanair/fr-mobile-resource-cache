package main

import (
	"encoding/json"
	"testing"
)

func mockSyncDocument() []byte {
	return []byte(`{
    "_rev": "1-aaaaa",
    "_prop": true,
    "key": "value"
    }`)
}

func TestCleanupSyncDocument(t *testing.T) {
	result, err := cleanupSyncDocument(mockSyncDocument())

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
