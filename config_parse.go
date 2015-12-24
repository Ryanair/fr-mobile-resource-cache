package main

import (
	"encoding/json"
	"io"
)

//DEBUG flag toggles log output
var DEBUG bool
var config Config
var authConfig Auth
var webPort = 8181

//DefaultAttachmentDoc is used as a template
//for attachment documents. If no attachment doc exists,
//this will be used instead
var DefaultAttachmentDoc string

//Config represents all resource and sync endpoints
type Config struct {
	ResourcesDir string `json:"resourcesDir"`
	SyncURL      string `json:"syncUrl"`
	Bucket       string `json:"bucket"`
	Auth         Auth   `json:"auth"`
}

//ConfigJSON represents global config properties
type ConfigJSON struct {
	Debug                bool   `json:"debug"`
	WebPort              int    `json:"webPort"`
	Resources            Config `json:"resources"`
	DefaultAttachmentDoc string `json:"default_attachment_doc"`
}

//Auth authentication config
type Auth struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	SimpleAuth bool   `json:"simpleAuth"`
	ServerURL  string `json:"serverUrl"`
}

//Export parses and serializes config.json
func (r ConfigJSON) Export() (Config, error) {
	result := Config{}
	result.ResourcesDir = r.Resources.ResourcesDir
	result.SyncURL = r.Resources.SyncURL
	result.Bucket = r.Resources.Bucket

	authConfig.Username = r.Resources.Auth.Username
	authConfig.Password = r.Resources.Auth.Password
	authConfig.SimpleAuth = r.Resources.Auth.SimpleAuth
	authConfig.ServerURL = r.Resources.Auth.ServerURL

	//set global DEBUG value
	DEBUG = r.Debug

	//set global endpoint values
	config = result

	//set default attachment document
	DefaultAttachmentDoc = r.DefaultAttachmentDoc

	webPort = r.WebPort

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
