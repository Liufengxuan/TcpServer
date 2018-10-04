package handle

import (
	"TcpServer/com"
	"TcpServer/models/dbops"
	"TcpServer/models/recode"
	"log"

	"github.com/astaxie/beego/orm"
)

/**************************验证并解析身份**********************************/
func addUser(curInfo *context) {
	curInfo.conn.Write([]byte("创建成功"))
}

/**************************验证并解析身份**********************************/
func auth(curInfo *context) { //身份验证成功返回02；程序出错返回4500；命令错误返回4104；数据库查询出错4001
	var rst string
	var serverErr string
	defer func() {
		if serverErr != "" {
			curInfo.endFlag = true
			curInfo.conn.Write([]byte(serverErr))
		} else if rst != "" {
			curInfo.conn.Write([]byte(rst))
			log.Printf("[%s:登陆成功]\n", curInfo.sessionInfo.UserIp)
		}

	}()
	//--------------------------------------------------
	if len(curInfo.cmds) != 3 {
		serverErr = recode.RECODE_LOGINPARAMERR
		return
	}
	if curInfo.cmds[0] == "LOGIN" {
		curInfo.userInfo.Name = curInfo.cmds[1]
		err := dbops.DBContext.Read(&curInfo.userInfo, "Name")
		if err != nil {
			if err == orm.ErrNoRows {
				log.Println("查询用户数据出错", err)
				serverErr = recode.RECODE_USERERR
				return
			}
			serverErr = recode.RECODE_DBERR
			return
		}

		pwd, err := com.GetMD5(curInfo.cmds[2])
		if err != nil {
			log.Println("MD5生成出错", err)
			serverErr = recode.RECODE_SERVERERR
			return
		}
		if pwd == curInfo.userInfo.PassWord {
			uuid, err := com.GetUUID()
			if err != nil {
				log.Println("uuid生成出错", err)
				serverErr = recode.RECODE_SERVERERR
				return
			}
			curInfo.sessionInfo.UserIp = curInfo.conn.RemoteAddr().String()
			curInfo.sessionInfo.SId = uuid
			rst = recode.RECODE_LOGINOK
			return
		}

		serverErr = recode.RECODE_PWDERR
		return
	}

}

/**************************未识别的命令**********************************/
func unIdentified(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_UNIDENTIFIED))
}

/**************************拒绝连接。**********************************/
func turnDown(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_PARAMERR))
	curInfo.endFlag = true
}

/**************************退出指令**********************************/
func exit(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_OK))
	curInfo.endFlag = true
}
