package main

import (
	"github.com/evanphx/json-patch"
)

func compare(a, b []byte) bool {
	result := jsonpatch.Equal(a, b)

	return result
}
