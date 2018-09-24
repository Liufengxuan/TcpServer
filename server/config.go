package main

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/config"
)

func init() {
	//获取配置文件信息  begin
	log.Println("[获取配置文件信息]")
	conf, err := config.NewConfig("ini", "./server.conf")
	if err != nil {
		fmt.Println("read ./config.conf err :", err)
		return
	}
	ip = conf.String("server::listen_ip")
	port = conf.String("server::listen_port")
	log.Println("[以获取到配置文件信息]")
	//获取配置文件信息  end
}
