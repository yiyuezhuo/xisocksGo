package main

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
)

type ParseHeaderResult struct {
	Method string
	Host   string
	Port   string
}

/*
func Parse(header []byte) ParseHeaderResult {
	Method := "GET" // "GET", "POST", "CONNECT"
	Host := "127.0.0.1"
	Port := 80
	return ParseHeaderResult{Method, Host, Port}
}
*/
func Parse(header []byte) (*ParseHeaderResult, error) {
	//https://stackoverflow.com/questions/33963467/parse-http-requests-and-responses-from-text-file-in-go
	//https://golang.org/pkg/net/http/#ReadRequest
	//https://golang.org/pkg/net/http/#Request
	reader := bytes.NewReader(header)
	bufio_reader := bufio.NewReader(reader)
	req, err := http.ReadRequest(bufio_reader)
	if err != nil {
		return nil, err
	}

	res := new(ParseHeaderResult)
	res.Method = req.Method
	if strings.Contains(req.URL.Host, ":") {
		host_port := strings.Split(req.URL.Host, ":")
		res.Host = host_port[0]
		res.Port = host_port[1]
	} else {
		res.Host = req.URL.Host
		res.Port = "80" // Is it possible that port take 443 in some situation?
	}
	return res, nil
}
