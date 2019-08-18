// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/yiyuezhuo/xisocksGo/fakewebsite"
	//"github.com/yiyuezhuo/cdnsocksGo/fakewebsite"
)

var upgrader = websocket.Upgrader{} // upgrader := websocket.Upgrader{} can't be used outside a function

func websocket_to_socket(websocket_conn *websocket.Conn, target_conn net.Conn) {
	for {
		_, message, err := websocket_conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Println("pipe:", len(message), message[:10])

		target_conn.Write(message)
	}
}
func socket_to_websocket(websocket_conn *websocket.Conn, target_conn net.Conn) {
	buf := make([]byte, 8192)
	for {
		count, err := target_conn.Read(buf)
		if err != nil {
			log.Println("Error Server read from target server fail", err)
			return
		}

		err = websocket_conn.WriteMessage(websocket.BinaryMessage, buf[:count])
		if err != nil {
			log.Println("Error websocket_conn.WriteMessage:", err)
			break
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	// Keep echo to help debugging
	fmt.Println("A fake connection have been tried")
	c, err := upgrader.Upgrade(w, r, nil) // upgrade from http to websocket connection(c)
	if err != nil {
		log.Println("Error upgrade:", err)
		return
	}
	defer c.Close()

	fmt.Println("The http connection have been upgraded to Websocket")

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("Error c.ReadMessage:", err)
		return
	}

	//fmt.Println("message[0]:", message[0])

	if message[0] == 5 { // socks 5
		fmt.Println("Assuming socks5")
		update_socks5_proxy(c, mt, message)
	} else {
		fmt.Println("Assuming http")
		update_http_proxy(c, mt, message)
	}
}

func update_http_proxy(c *websocket.Conn, mt int, message []byte) {
	log.Println("First read:", len(message), string(message[:20]))
	header, err := Parse(message)
	if err != nil {
		log.Println("Error Parse header:", header)
		return
	}

	host_port := header.Host + ":" + header.Port
	target_conn, err := net.Dial("tcp", host_port)
	if err != nil {
		log.Println("Error Dial to target conn:", host_port)
		return
	}

	if header.Method == "CONNECT" {
		err = c.WriteMessage(mt, []byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			log.Println("Error response to connection established:", err)
			return
		}
	} else {
		target_conn.Write(message)
	}

	go socket_to_websocket(c, target_conn)
	websocket_to_socket(c, target_conn)
}

func update_socks5_proxy(c *websocket.Conn, mt int, message []byte) {
	buf := make([]byte, 2048)
	buf_dummy := make([]byte, 0)
	copy(buf, message) // Assuming is small enough so buf can store all of its content
	wca := &WsConnAdaptor{c, buf_dummy}
	//remote_host, err := socks5_handshake(wca, message, len(message))
	remote_host, err := socks5_handshake(wca, buf, len(message))
	if err != nil {
		fmt.Println("SOCKS5 handshake fail", err)
		return
	}
	target_conn, err := net.Dial("tcp", remote_host)
	if err != nil {
		log.Println("Error Dial to target conn:", remote_host)
		return
	}

	go socket_to_websocket(c, target_conn)
	websocket_to_socket(c, target_conn)
}

func echo(w http.ResponseWriter, r *http.Request) {
	// Keep echo to help debugging
	fmt.Println("A fake connection have been tried")
	c, err := upgrader.Upgrade(w, r, nil) // upgrade from http to websocket connection(c)
	if err != nil {
		log.Print("Error echo upgrade:", err)
		return
	}
	defer c.Close()

	fmt.Println("The http connection have been upgraded to Websocket")

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error echo read_message:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("Error echo write_message:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(fakewebsite.Sadpanda)
	w.Write([]byte(fakewebsite.Sadpanda))
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	fmt.Println("start xisocks 0.4 server")

	//addr := "127.0.0.1:80"
	config := loadConfig()
	addr := config.ListenIp + ":" + strconv.Itoa(config.ListenPort)

	http.HandleFunc("/upload", upload) // fake url to cheat GFW reviewer
	http.HandleFunc("/echo", echo)     // fake url to cheat GFW reviewer
	http.HandleFunc("/", home)
	http.HandleFunc("/sadpanda.jpg", fakewebsite.SendSadPanda)

	fmt.Println(addr, upgrader)
	//var upgrader := websocket.Upgrader{}
	if !config.TLS {
		fmt.Println("launch without TLS")
		log.Fatal(http.ListenAndServe(addr, nil))
	} else {
		fmt.Println("launch with TLS")
		fmt.Println("crt:", config.Crt, "key", config.Key)
		log.Fatal(http.ListenAndServeTLS(addr, config.Crt, config.Key, nil))
	}
	fmt.Print("wtf")
}
