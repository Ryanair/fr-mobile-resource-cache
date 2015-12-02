package main

import (
	"github.com/evanphx/json-patch"
)

func compare(a, b []byte) bool {
	result := jsonpatch.Equal(a, b)

	return result
}

func diff(a, b []byte) ([]byte, error) {
	return jsonpatch.CreateMergePatch(a, b)
}
