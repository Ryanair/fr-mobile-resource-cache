package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/couchbaselabs/logg"
	"github.com/go-fsnotify/fsnotify"
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
	startWebServer()

	resources, _ := getLocalResources()
	patchFiles(resources)
	trackScanFiles(len(resources))

	dirList, err := getDirectories()
	if err != nil {
		logg.LogPanic("Error scanning directories", err)
	}

	newFolderWatcher(dirList)
}

func startWebServer() {
	sock, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", webPort))
	if err != nil {
		logg.LogPanic("Error starting web server : %v", err)
	}
	go func() {
		fmt.Println("HTTP now available at port ", webPort)
		http.Serve(sock, nil)
	}()

}

func newFolderWatcher(dirList []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				logg.LogTo(TagLog, "New Event %v", event)
				//rename reports the old filename
				if event.Op&fsnotify.Remove != fsnotify.Remove && event.Op&fsnotify.Rename != fsnotify.Rename {
					// localResource, _ := NewLocalDocument(event.Name)
					// patchFiles([]LocalResource{localResource})
				}
				// TODO: handle deletes
			case err := <-watcher.Errors:
				logg.LogTo(TagError, "%v", err)
			}
		}
	}()

	for _, dir := range dirList {
		logg.LogTo(TagLog, "attaching watcher to %s", dir)
		err = watcher.Add(dir)
		if err != nil {
			logg.LogPanic("Error attaching fs watcher : %v", err)
		}
	}
	<-done
}
