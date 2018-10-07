package handle

import (
	"TcpServer/com"
	"TcpServer/loging"
	"TcpServer/models/dbops"
	"TcpServer/models/recode"
	"net"

	"github.com/astaxie/beego/orm"
)

type context struct {
	userInfo    dbops.User
	sessionInfo dbops.Session
	//当前用户的   连接对象
	conn net.Conn
	//当前用户的    命令
	cmds []string

	cmdLength int
	//当前用户的    会话结束标志
	endFlag bool
}

/**************************验证并解析身份**********************************/
func (ctxt *context) addUser() { //2009权限不足、2007用户信息重复、2002命令格式不对、4500服务器出错、4005添加数据出错、1001成功
	var respCode string
	var err error
	defer func() {
		respByte(ctxt, respCode)
	}()
	//判断是不是管理员。
	if !isAdmin(ctxt.userInfo.Name) {
		respCode = recode.RECODE_NOPERMISSION
		return
	}

	user := dbops.User{Name: ctxt.cmds[2]}
	err = dbops.DBContext.Read(&user, "Name")
	if err != nil && err != orm.ErrNoRows {
		respCode = recode.RECODE_DBERR
		return
	}
	if user.Id != 0 {
		respCode = recode.RECODE_REPEATUSER
		return
	}
	if ctxt.cmdLength != 4 {
		respCode = recode.RECODE_CMDFORMATERR
		return
	}

	//通过验证后、给用户创建数据
	newUser := new(dbops.User)
	newHomeDir := new(dbops.FileInfo)
	newUser.PassWord, err = com.GetMD5(ctxt.cmds[3])
	if err != nil {
		respCode = recode.RECODE_SERVERERR
		return
	}
	newUser.Name = ctxt.cmds[2]
	newHomeDir.IsDir = true
	newHomeDir.IsHome = true
	newHomeDir.Name = ctxt.cmds[2] + "Home"
	//开启一个事务
	err = dbops.DBContext.Begin()
	_, err = dbops.DBContext.Insert(newUser)
	newHomeDir.User = newUser
	_, err = dbops.DBContext.Insert(newHomeDir)
	if err != nil {
		rollBackErr := dbops.DBContext.Rollback()
		if rollBackErr != nil {
			loging.Loger.Error("[添加用户数据回滚失败、可能产生了脏数据 err=%s]\n", rollBackErr)
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
func (ctxt *context) auth() { //身份验证成功返回1002；程序出错返回4500；命令错误返回2004；数据库查询出错2008
	var rst string
	var serverErr string
	defer func() {
		if serverErr != "" {
			ctxt.endFlag = true
			ctxt.conn.Write([]byte(serverErr))
		} else if rst != "" {
			ctxt.conn.Write([]byte(rst))
			loging.Loger.Info("[%s:登陆成功]\n", ctxt.sessionInfo.UserIp)
		}
	}()
	//--------------------------------------------------
	if len(ctxt.cmds) != 3 {
		serverErr = recode.RECODE_LOGINPARAMERR
		return
	}
	if ctxt.cmds[0] == "LOGIN" {
		ctxt.userInfo.Name = ctxt.cmds[1]
		err := dbops.DBContext.Read(&ctxt.userInfo, "Name")
		if err != nil {
			if err == orm.ErrNoRows {
				loging.Loger.Error("查询用户数据出错", err)
				serverErr = recode.RECODE_USERERR
				return
			}
			serverErr = recode.RECODE_DBERR
			return
		}

		pwd, err := com.GetMD5(ctxt.cmds[2])
		if err != nil {
			loging.Loger.Error("MD5生成出错", err)
			serverErr = recode.RECODE_SERVERERR
			return
		}
		if pwd == ctxt.userInfo.PassWord {
			uuid, err := com.GetUUID()
			if err != nil {
				loging.Loger.Error("uuid生成出错", err)
				serverErr = recode.RECODE_SERVERERR
				return
			}
			ctxt.sessionInfo.UserIp = ctxt.conn.RemoteAddr().String()
			ctxt.sessionInfo.SId = uuid
			rst = recode.RECODE_LOGINOK
			return
		}

		serverErr = recode.RECODE_PWDERR
		return
	}

}

/**************************未识别的命令**********************************/
func (ctxt *context) unIdentified() {
	ctxt.conn.Write([]byte(recode.RECODE_UNIDENTIFIED))
}

/**************************拒绝连接。**********************************/
func (ctxt *context) turnDown() {
	ctxt.endFlag = true
	ctxt.conn.Write([]byte(recode.RECODE_PARAMERR))

}

/**************************退出指令**********************************/
func (ctxt *context) exit() {
	ctxt.endFlag = true
	ctxt.conn.Write([]byte(recode.RECODE_OK))
}
