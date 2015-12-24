package main

import "expvar"

var (
	incrementFiles int64
	watchfiles     = expvar.NewInt("watchfiles")
)

func trackScanFiles(i int) {
	incrementFiles += int64(i)
	watchfiles.Set(incrementFiles)
}
