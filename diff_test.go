package main

import (
	"testing"
)

func TestCompare(t *testing.T) {
	doc := []byte(`{
    "key" : "value"
  }`)

	if compare(doc, doc) == false {
		t.Errorf("compare assertion failed")
	}

}

func TestDiff(t *testing.T) {
	doc1 := []byte(`{
  	"key" : "value"
  }`)

	doc2 := []byte(`{
    "key" : "value1"
  }`)

	patch, err := diff(doc1, doc2)

	if err != nil {
		t.Errorf("Error generating patch", err)
	}

	str := string(patch)
	if str != `{"key":"value1"}` {
		t.Error("diff assertion failed")
	}
}
