package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config contiains configuration
type Config struct {
	ConfigFile    	string
	APIToken      	string 	`json:"token"`
	BaseURL       	string 	`json:"base-url"`
	Email         	string 	`json:"email"`
	ProjectID     	string 	`json:"projectID"`
	EnvID         	string 	`json:"envID"`
	ReferenceType 	string 	`json:"defaultReferenceType"`
	Debug			bool 	`json:"debug"`
}

//Init new config
func (c *Config) Init() {
	configFile, err := ioutil.ReadFile(c.ConfigFile)
	check(err)
	if err = json.Unmarshal(configFile, c); err != nil {
		panic(err)
	}
}
