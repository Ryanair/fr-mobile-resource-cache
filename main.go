package main

import (
	"bufio"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/couchbaselabs/logg"
)

var (
	endpoint              string
	configFileDescription = "The name of the config file.  Defaults to 'config.json'"
	configFileName        = kingpin.Arg("config file name", configFileDescription).Default("config.json").String()
)

func main() {
	kingpin.Parse()
	if *configFileName == "" {
		kingpin.Errorf("Config file name missing")
		return
	}
	configFile, err := os.Open(*configFileName)
	if err != nil {
		logg.LogPanic("Unable to open file: %v.  Err: %v", *configFileName, err.Error())
		return
	}
	defer configFile.Close()

	configReader := bufio.NewReader(configFile)

	parseConfigFile(configReader)

	files := scanResourcesDir()

	for _, file := range files {
		localDocument, fileName := readFileContents(file)

		syncDocument, _ := getDocument(fileName)

		syncDocument = cleanupSyncDocument(syncDocument)

		//if the document does not exist, or there is a difference in revisions, post it
		if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localDocument, syncDocument) {
			postDocument(localDocument, fileName)
		}
	}
}
