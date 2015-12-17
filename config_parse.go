package main

import (
	"encoding/json"
	"io"
)

//DEBUG flag toggles log output
var DEBUG bool
var config Config

//DefaultAttachmentDoc is used as a template
//for attachment documents. If no attachment doc exists,
//this will be used instead
var DefaultAttachmentDoc string

//Config represents all resource and sync endpoints
type Config struct {
	ResourcesDir string `json:"resourcesDir"`
	SyncURL      string `json:"syncUrl"`
	Bucket       string `json:"bucket"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

//ConfigJSON represents global config properties
type ConfigJSON struct {
	Debug                bool   `json:"debug"`
	Resources            Config `json:"resources"`
	DefaultAttachmentDoc string `json:"default_attachment_doc"`
}

//Export parses and serializes config.json
func (r ConfigJSON) Export() (Config, error) {
	result := Config{}
	result.ResourcesDir = r.Resources.ResourcesDir
	result.SyncURL = r.Resources.SyncURL
	result.Bucket = r.Resources.Bucket
	result.Username = r.Resources.Username
	result.Password = r.Resources.Password

	//set global DEBUG value
	DEBUG = r.Debug

	//set global endpoint values
	config = result

	//set default attachment document
	DefaultAttachmentDoc = r.DefaultAttachmentDoc

	return result, nil
}

func parseConfigFile(r io.Reader) error {
	configJSON := ConfigJSON{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&configJSON); err != nil {
		return err
	}
	_, err := configJSON.Export()

	return err
}
