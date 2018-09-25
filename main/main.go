package main

import (
	"log"
	"net"
)

var ip string
var port string

func main() {
	//监听 begin
	address := ip + ":" + port
	listener, err := net.Listen("tcp", address)
	defer listener.Close()
	if err != nil {
		log.Printf("[listener err=%s]\n", err)
		return
	}
	log.Printf("[准备监听 %s]\n", address)

	listenerErr := waitConnection(listener)
	if err != nil {
		log.Printf("[监听程序出现异常：err=%s]\n", listenerErr)
		return
	}

}
