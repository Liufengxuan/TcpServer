package loging

import (
	"TcpServer/com"
	"log"
	"strconv"

	"github.com/astaxie/beego/logs"
)

var Loger *logs.BeeLogger

func init() {
	initLogs()
}
func initLogs() {
	var debug bool
	var maxlines int
	var maxdays int
	var filename string
	var err1 error
	Loger = logs.NewLogger()

	cfg, err := com.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	debug, err1 = strconv.ParseBool(cfg.GetString("logs::Debug"))
	if err1 != nil {
		debug = false
		log.Println("配置文件logs::Debug 读取失败、或读取失败；debug模式默认已被设为关闭")
	}

	if debug {
		name := cfg.GetString("logs::log_filename")
		if name == "" {
			name = "ser.log"
			log.Println("[log_filename获取失败，已设为 'ser.log']")
		}
		filename = name

		maxlines, err1 = cfg.GetInt("logs::log_maxlines")
		if err1 != nil || maxlines < 50 {
			maxlines = 150
			log.Println("log_maxlines获取失败或小于50行，已设为 ‘150’")
		}

		maxdays, err1 = cfg.GetInt("logs::log_maxdays")
		if err1 != nil || maxdays < 1 {
			maxdays = 7
			log.Println("log_maxdays获取失败或小于1，已经设为 ‘7’")
		}

		var config = `{"filename":"` +
			filename + `","maxlines":` +
			strconv.Itoa(maxlines) + `,"maxdays":` +
			strconv.Itoa(maxdays) + `,"separate":["error","warning"]}`
		err2 := Loger.SetLogger(logs.AdapterMultiFile, config)
		if err2 != nil {
			log.Println("日志模块配置失败")
		}
	} else {
		Loger.SetLogger(logs.AdapterConsole)
	}

}
