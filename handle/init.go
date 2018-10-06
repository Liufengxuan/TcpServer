package handle

import (
	"TcpServer/com"
	"TcpServer/loging"
	"log"
	"runtime"
	"strings"
	"time"
)

//缓冲区大小
var readBufferSize int
var writeBufferSize int

//当前系统的换行符长度。
var lineBreakLength int

//当前系统的管理员账户
var adminList []string

//等待用户消息超时时间
var respTimeout time.Duration

func init() {
	initSysType()
	initBufferSize()
	initAdminUserInfo()
	initRespTimeout()
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
		loging.Loger.Error("<缓冲区配置读取失败、已经设置为 readBuffer=4096,writeBuffer=2048>")
	}
}

func initAdminUserInfo() {
	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
	}

	str := cfg.GetString("User::admin_userNmae")
	if str == "" {
		log.Println("[未从配置文件读取到管理员名称]")
		loging.Loger.Error("[未从配置文件读取到管理员名称]")
		return
	}
	adminList = strings.Split(str, "*")

}
func initRespTimeout() {
	var err1 error
	var second int
	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
	}

	second, err1 = cfg.GetInt("session::responseTimeout")
	if err1 != nil {
		respTimeout = time.Second * 90
		log.Println("<responseTimeout配置项读取失败、已经设置默认值 responseTimeout=90>")
	}
	respTimeout = time.Second * time.Duration(second)

	// switch second {
	// case 8:
	// 	respTimeout = time.Second * 8
	// case 16:
	// 	respTimeout = time.Second * 16
	// case 32:
	// 	respTimeout = time.Second * 32
	// case 64:
	// 	respTimeout = time.Second * 64
	// case 128:
	// 	respTimeout = time.Second * 128
	// case 256:
	// 	respTimeout = time.Second * 256
	// case 512:
	// 	respTimeout = time.Second * 512
	// case 1024:
	// 	respTimeout = time.Second * 1024
	// default:
	// 	respTimeout = time.Second * 90
	// 	log.Println("<responseTimeout配置项值无效、已经设置默认值 responseTimeout=90>")
	// }

}
