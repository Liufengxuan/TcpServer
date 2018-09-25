package main

import (
	"log"
	"net"
)

func waitConnection(listener net.Listener) error {
	log.Println("[等待用户连接.......]")
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		handelConnection(conn)
	}
	//监听 end
}

func handelConnection(conn net.Conn) {
	log.Printf("[%s连接]", conn.RemoteAddr())
}
