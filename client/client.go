package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func main() {
	config := loadConfig()
	RemoteIp := config.RemoteIp
	RemotePort := strconv.Itoa(config.RemotePort)
	LocalIp := config.LocalIp
	LocalPort := strconv.Itoa(config.LocalPort)
	TLS := config.TLS

	if TLS {
		fmt.Println("Use TLS")
	}

	// https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go
	l, err := net.Listen("tcp", LocalIp+":"+LocalPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + LocalIp + ":" + LocalPort)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, RemoteIp+":"+RemotePort, TLS)
	}
}

// Handles incoming requests.
func handleRequest(local_c net.Conn, addr string, TLS bool) {
	// Make a buffer to hold incoming data.
	defer local_c.Close()

	var u url.URL

	if !TLS {
		u = url.URL{Scheme: "ws", Host: addr, Path: "/upload"}
	} else {
		u = url.URL{Scheme: "wss", Host: addr, Path: "/upload"}
	}
	log.Printf("connecting to %s", u.String())

	fmt.Println("before dial")
	remote_c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	fmt.Println("after dial")
	if err != nil {
		log.Println("Error Dial:", err)
		return
	}
	defer remote_c.Close()

	//done := make(chan struct{})

	go func() {
		//defer close(done)
		for {
			_, message, err := remote_c.ReadMessage()
			if err != nil {
				log.Println("Error remore_c.ReadMessage:", err)
				return
			}
			//log.Printf("recv: %s", message)
			fmt.Println("remote to local:", len(message), message[:10])

			local_c.Write(message)
		}
	}()

	buf := make([]byte, 8192)

	for {
		readLen, err := local_c.Read(buf)
		if err != nil {
			log.Println("Error local_c.read:", err)
			return
		}

		fmt.Println("local to remote:", readLen, buf[:10])
		remote_c.WriteMessage(websocket.BinaryMessage, buf[:readLen])
	}

}
