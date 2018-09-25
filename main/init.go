package main

import (
	"TcpServer/com"
	"log"
)

func init() {
	//获取配置文件信息  begin
	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	ip = cfg.GetString("server::listen_ip")
	port = cfg.GetString("server::listen_port")
	//获取配置文件信息  end
}
