package main

import "expvar"

var (
	watchfiles = expvar.NewInt("watchfiles")
)

func trackScanFiles(i int) {
	watchfiles.Set(int64(i))
}
