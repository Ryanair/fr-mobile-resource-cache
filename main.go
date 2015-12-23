package main

import (
	"bufio"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/couchbaselabs/logg"
)

var (
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
	resources, _ := getLocalResources()
	for _, resource := range resources {
		logg.LogTo(TagLog, "DocumentID : %s ; FileName : %s ; Attachment : %s ; Content : %s", resource.ResourceID, resource.FileName, resource.Attachment, string(resource.Content))
	}
	//first time scan of the directories
	// files, err := getLocalResources()
	//
	// if err != nil {
	// 	logg.LogPanic("Error scanning resource directory: %v", err)
	// }
	//
	// patch, err := patchFiles(files)
	//
	// if err != nil {
	// 	logg.LogTo(TagError, "%v", err)
	// }
	//
	// if len(patch) > 0 {
	// 	logg.LogTo(TagDiff, "%s", patch)
	// }

	//register the fs watcher
	// TODO: refactor to get rid of the double recursion in the fist run of the app
	// dirList, err := getDirectories()
	// if err != nil {
	// 	logg.LogPanic("Error scanning directories : %v", err)
	// }
	//
	// newFolderWatcher(dirList)
}

func newFolderWatcher(dirList []string) {
	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer watcher.Close()
	//
	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-watcher.Events:
	// 			logg.LogTo(TagLog, "New Event %v", event)
	// 			//rename reports the old filename
	// 			if event.Op&fsnotify.Remove != fsnotify.Remove && event.Op&fsnotify.Rename != fsnotify.Rename {
	// 				localResource, _ := NewLocalDocument(event.Name)
	// 				patchFiles([]LocalResource{localResource})
	// 			}
	// 			// TODO: handle deletes
	// 		case err := <-watcher.Errors:
	// 			logg.LogTo(TagError, "%v", err)
	// 		}
	// 	}
	// }()
	//
	// for _, dir := range dirList {
	// 	logg.LogTo(TagLog, "attaching watcher to %s", dir)
	// 	err = watcher.Add(dir)
	// 	if err != nil {
	// 		logg.LogPanic("Error attaching fs watcher : %v", err)
	// 	}
	// }
	// <-done
}
