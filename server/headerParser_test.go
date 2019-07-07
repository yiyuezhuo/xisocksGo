package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT
	connect_header := []byte("CONNECT server.example.com:80 HTTP/1.1\r\nHost: server.example.com:80\r\nProxy-Authorization: basic aGVsbG86d29ybGQ=\r\n\r\n")
	res, err := Parse(connect_header)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Method:", res.Method, "Host:", res.Host, "Port:", res.Port)
	get_header := []byte("GET /tutorials/other/top-20-mysql-best-practices/ HTTP/1.1\r\nHost: net.tutsplus.com\r\nUser-Agent: Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.5) Gecko/20091102 Firefox/3.5.5 (.NET CLR 3.5.30729)\r\n\r\n")
	res, err = Parse(get_header)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Method:", res.Method, "Host:", res.Host, "Port:", res.Port)
}
