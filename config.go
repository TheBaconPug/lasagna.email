package main

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigStruct struct {
	Port     string   `json:"port"`
	Domains  []string `json:"domains"`
	MongoURI string   `json:"mongoURI"`
}

var Config ConfigStruct

func LoadConfig() {
	file, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(file, &Config)
}
