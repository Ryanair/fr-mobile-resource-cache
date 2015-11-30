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
	//set logging
	logg.LogKeys[TagError] = true
	logg.LogKeys[TagLog] = true

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

	err = parseConfigFile(configReader)
	if err != nil {
		logg.LogPanic("Erro parsing the config file: %v", err)
	}

	files, err := scanResourcesDir()

	if err != nil {
		logg.LogPanic("Error scanning resource directory: %v", err)
	}

	for _, file := range files {
		localDocument, fileName, err := readFileContents(file)

		if err != nil {
			// log the error, don't stop execution if a file fails to read
			logg.LogTo(TagError, "Error reading file contents: %v", err)
			continue
		}

		syncDocument, _, err := getDocument(fileName)
		if err != nil {
			logg.LogTo(TagError, "Error reading sync document: %v", err)
			continue
		}

		syncDocument, err = cleanupSyncDocument(syncDocument)
		if err != nil {
			logg.LogTo(TagError, "Error cleaning up sync document: %v", err)
			continue
		}
		//if the document does not exist, or there is a difference in revisions, post it
		if (len(syncDocument) == 0 || syncDocument == nil) || !compare(localDocument, syncDocument) {
			err := postDocument(localDocument, fileName)
			if err != nil {
				logg.LogTo(TagError, "Error saving document: %v", err)
				continue
			}
		}
	}
}
