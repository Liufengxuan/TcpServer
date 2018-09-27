package handle

import (
	"TcpServer/com"
	"log"
	"runtime"
)

//缓冲区大小
var readBufferSize int
var writeBufferSize int

//当前系统的换行符长度。
var lineBreakLength int

func init() {
	initSysType()
	initBufferSize()
}

func initSysType() {
	currentSys := runtime.GOOS
	if currentSys == "windows" {
		lineBreakLength = 2
	} else {
		lineBreakLength = 1
	}
	//log.Printf("[当前系统为： %s ]\n", currentSys)
}

func initBufferSize() {
	var err1 error
	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
	}

	readBufferSize, err1 = cfg.GetInt("socket::readBuf_size")

	writeBufferSize, err1 = cfg.GetInt("socket::writeBuf_size")
	if err1 != nil {
		readBufferSize = 4096
		writeBufferSize = 2048
		log.Println("<缓冲区配置读取失败、已经设置为 readBuffer=4096,writeBuffer=2048>")
	}
	log.Println("[已初始化缓冲区设置]")
}
