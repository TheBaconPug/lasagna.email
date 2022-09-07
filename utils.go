package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

type Config struct {
	Port          string `json:"port"`
	Domain        string `json:"domain"`
	MailgunAPIKey string `json:"mailgun_api_key"`
}

func LoadConfig() Config {
	var config Config
	file, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(file, &config)
	return config
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
