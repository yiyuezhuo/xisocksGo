package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ListenIp   string
	ListenPort int
	TLS        bool
	Crt        string
	Key        string
}

func loadConfig() Config {
	// loaf config from json
	jsonFile, err := os.Open("config-server.json")
	if err != nil {
		log.Panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("ListenIp:", config.ListenIp, "ListenPort:", config.ListenPort, "TLS:", config.TLS,
		"crt:", config.Crt, "key:", config.Key)
	return config
}
