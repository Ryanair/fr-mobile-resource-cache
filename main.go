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

func init() {
	//set logging
	logg.LogKeys[TagError] = true
	logg.LogKeys[TagDiff] = true

	if DEBUG == false {
		logg.LogKeys[TagLog] = true
	}

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
}

func main() {
	files, err := scanResourcesDir()

	if err != nil {
		logg.LogPanic("Error scanning resource directory: %v", err)
	}

	patch, err := patchFiles(files)

	if err != nil {
		logg.LogTo(TagError, "%v", err)
	}

	if len(patch) > 0 {
		logg.LogTo(TagDiff, "%s", patch)
	}
}
