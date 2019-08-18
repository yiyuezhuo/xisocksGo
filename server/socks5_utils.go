package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type ReadWriteCloseBoy interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}

type WsConnAdaptor struct {
	conn *websocket.Conn
	buf  []byte
}

func (wca *WsConnAdaptor) Write(b []byte) (n int, err error) {
	//fmt.Println("wca WriteMessage:", b)
	err = wca.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), err
}

func (wca *WsConnAdaptor) Read(b []byte) (n int, err error) {
	//fmt.Println("len of buf:", len(wca.buf), "cap:", cap(wca.buf))
	if len(wca.buf) > 0 {
		readLen := copy(b, wca.buf)
		wca.buf = wca.buf[readLen:]
		return readLen, nil
	} else {
		_, message, err := wca.conn.ReadMessage()
		//fmt.Println("wca ReadMessage:", message, err)
		if err != nil {
			fmt.Println("wca read error", err)
			return 0, err
		}
		readLen := copy(b, message)
		wca.buf = message[readLen:]
		//fmt.Println("wca readLen:", readLen, "len(wca.buf)", len(wca.buf))
		return readLen, nil
	}
}
func (wca *WsConnAdaptor) Close() (err error) {
	err = wca.conn.Close()
	return err
}

func socks5_handshake(local_c ReadWriteCloseBoy, buf []byte, readLen int) (remote_host string, err error) {
	//buf := make([]byte, 2048)

	//readLen, err := local_c.Read(buf)
	if readLen == 0 {
		_, err = local_c.Read(buf)
	}

	if err != nil {
		log.Println("Error local_c.read:", err)
		return "", err
	}

	//fmt.Println("local to remote:", readLen, buf[:10])

	version := buf[0]
	nmethods := buf[1]
	methods := buf[2 : 2+nmethods]

	//fmt.Println("version:", version, "nmethods:", nmethods, "methods:", methods)
	// version: 5 nmethods: 1 methods: [0]
	if version == 4 {
		fmt.Println("Socks4 protolcol, it support only Socks5")
		return "", fmt.Errorf("Socks4 protolcol, it support only Socks5")
	} else if version != 5 {
		fmt.Println("Unknow protocol detected, it support only socks5")
		return "", fmt.Errorf("Unknow protocol detected, it support only socks5")
	}

	/*
		if (nmethods != 1) || (methods[0] != 0) {
			fmt.Println("Trying to use complex authentication, it support only NO AUTHENTICATION REQUIRED")
			return "", fmt.Errorf("Trying to use complex authentication, it support only NO AUTHENTICATION")
		}
	*/
	can_use_no_authentication := false
	for i := 0; i < int(nmethods); i++ {
		if methods[i] == 0 {
			can_use_no_authentication = true
			break
		}
	}
	if !can_use_no_authentication {
		fmt.Println("It support only NO AUTHENTICATION")
		return "", fmt.Errorf("It support only NO AUTHENTICATION")
	}

	buf[0] = 5 // socks5
	buf[1] = 0 // NO AUTHENTICATION REQUIRED
	_, err = local_c.Write(buf[:2])
	if err != nil {
		return "", err
	}

	//readLen, err = local_c.Read(buf)
	readLen, err = local_c.Read(buf)
	if err != nil {
		return "", err
	}

	//fmt.Println("readLen:", readLen, "content:", buf[:readLen])
	/*
		AUTHENTICATION, such as username/password is skiped
		Now we take the connect command:
			version	cmd		rsv		atyp	dst.addr		dst.port
			1 byte	1 byte	1 byte	1 byte	4 to 255 bytes	2 bytes
	*/
	version = buf[0]
	cmd := buf[1]
	rsv := buf[2]
	atyp := buf[3]

	//fmt.Println("version:", version, "cmd:", cmd, "rsv:", rsv, "atyp:", atyp)

	if cmd != 1 {
		fmt.Println("Support only CONNECT cmd")
		return "", fmt.Errorf("Support only CONNECT cmd")
	}
	if rsv != 0 {
		fmt.Println("Unknown reserve handling", rsv)
		return "", fmt.Errorf("Unknown reserve handling %d", rsv)
	}
	if atyp == 4 {
		fmt.Println("It can not handle ipv6")
		return "", fmt.Errorf("It can not handle ipv6")
	}

	//var remote_host string

	if atyp == 3 { // DOMAINNAME
		dst_addr_len := buf[4]
		dst_addr := buf[5 : 5+dst_addr_len]
		port := buf[5+dst_addr_len : 5+dst_addr_len+2]
		//fmt.Println("dst_addr_len:", dst_addr_len, " dst_addr:", dst_addr, " port:", port)
		remote_host = string(dst_addr) + ":" + strconv.Itoa(int(port[0])*256+int(port[1]))
		//fmt.Println("Parsed dst_addr:", string(dst_addr), " port:", int(port[0])*256+int(port[1]))
	} else if atyp == 1 { //IP V4
		dst_addr := buf[4 : 4+4] // or [4:4+4] ? // I can't construct a pure ipv4 request from my browser to verify
		port := buf[4+4 : 4+4+2] // or [4+4:4+4+2] ?
		//fmt.Println("dst_addr:", dst_addr, "port:", port)
		pp := fmt.Sprintf("%d.%d.%d.%d", dst_addr[0], dst_addr[1], dst_addr[2], dst_addr[3])
		remote_host = pp + ":" + strconv.Itoa(int(port[0])*256+int(port[1]))
		//fmt.Println("Parsed dst_addr:", pp, " port:", int(port[0])*256+int(port[1]))
	} else {
		fmt.Println("Unknown ATYP value", atyp)
		return "", fmt.Errorf("Unknown AYTP value %d", atyp)
	}

	buf[0] = 5    // version=socks5
	buf[1] = 0    // rep=succeeded
	buf[2] = 0    // rsv=0
	buf[3] = atyp // atyp = atyp

	for i := 0; i < 6; i++ { // fill dummy value to bnd.addr,bnd.port
		buf[4+i] = 0
	}

	/*
		In socks5 specification:
		https://www.ietf.org/rfc/rfc1928.txt
		It seems like that we should return :10 and curl accept it.
		But browser such as Chrome and Firefox reject it for some reason.
		I don't know but leave the clumsy condition statement to handle both.
	*/
	if atyp == 3 {
		_, err = local_c.Write(buf[:7])
	} else if atyp == 1 {
		_, err = local_c.Write(buf[:10])
	}

	//_, err = local_c.Write(buf[:7]) // why 7??? 7 works on browser but nor curl
	//_, err = local_c.Write(buf[:10]) // works on curl but not browser??

	fmt.Println(buf[:10])

	if err != nil {
		return "", err
	}

	return remote_host, nil
}
