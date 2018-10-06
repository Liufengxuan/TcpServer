package handle

import (
	"TcpServer/com"
	"TcpServer/models/dbops"
	"TcpServer/models/recode"
	"log"

	"github.com/astaxie/beego/orm"
)

/**************************验证并解析身份**********************************/
func addUser(curInfo *context) { //2009权限不足、2007用户信息重复、2002命令格式不对、4500服务器出错、4005添加数据出错、1001成功
	var respCode string
	var err error
	defer func() {
		respByte(curInfo, respCode)
	}()
	//判断是不是管理员。
	if !isAdmin(curInfo.userInfo.Name) {
		respCode = recode.RECODE_NOPERMISSION
		return
	}

	user := dbops.User{Name: curInfo.cmds[2]}
	err = dbops.DBContext.Read(&user, "Name")
	if err != nil && err != orm.ErrNoRows {
		respCode = recode.RECODE_DBERR
		return
	}
	if user.Id != 0 {
		respCode = recode.RECODE_REPEATUSER
		return
	}
	if curInfo.cmdLength != 4 {
		respCode = recode.RECODE_CMDFORMATERR
		return
	}

	//通过验证后、给用户创建数据
	newUser := new(dbops.User)
	newHomeDir := new(dbops.FileInfo)
	newUser.PassWord, err = com.GetMD5(curInfo.cmds[3])
	if err != nil {
		respCode = recode.RECODE_SERVERERR
		return
	}
	newUser.Name = curInfo.cmds[2]
	newHomeDir.IsDir = true
	newHomeDir.IsHome = true
	newHomeDir.Name = curInfo.cmds[2] + "Home"
	//开启一个事务
	err = dbops.DBContext.Begin()
	_, err = dbops.DBContext.Insert(newUser)
	newHomeDir.User = newUser
	_, err = dbops.DBContext.Insert(newHomeDir)
	if err != nil {
		rollBackErr := dbops.DBContext.Rollback()
		if rollBackErr != nil {
			log.Printf("[添加用户数据回滚失败、可能产生了脏数据 err=%s]\n", rollBackErr)
		}
		respCode = recode.RECODE_ADDDATAERR
	} else {
		commitErr := dbops.DBContext.Commit()
		if commitErr != nil {
			respCode = recode.RECODE_ADDDATAERR
		}
		respCode = recode.RECODE_OK
	}

}

/**************************验证并解析身份**********************************/
func auth(curInfo *context) { //身份验证成功返回1002；程序出错返回4500；命令错误返回2004；数据库查询出错2008
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
	curInfo.endFlag = true
	curInfo.conn.Write([]byte(recode.RECODE_OK))
}
