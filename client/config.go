package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	RemoteIp   string
	RemotePort int
	LocalIp    string
	LocalPort  int
	TLS        bool
}

func loadConfig() Config {
	// loaf config from json
	jsonFile, err := os.Open("config-client.json")
	if err != nil {
		log.Panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Panic(err)
	}
	return config
}
