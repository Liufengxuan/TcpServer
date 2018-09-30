package main

import (
	"TcpServer/com"
	"log"
)

var ip string
var port string
var reListenNum int
var maxProcs int

func init() {
	//获取配置文件信息  begin
	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	ip = cfg.GetString("server::listen_ip")
	port = cfg.GetString("server::listen_port")

	reListenNum, err = cfg.GetInt("server::mian_restartnumber")
	if err != nil {
		log.Println("[主进程意外重启次数配置项读取失败、已经设置为 3 次]")
		reListenNum = 3
	}

	maxProcs, err = cfg.GetInt("server::maxProcs")
	if err != nil {
		log.Println("[核心数配置项 读取失败、已经设置为 1核心]")
		maxProcs = 1
	}

	//获取配置文件信息  end
}
