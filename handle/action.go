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
	userInfo      dbops.User
	sessionInfo   dbops.Session
	currentFolder dbops.FileInfo
	//当前用户的连接对象
	conn net.Conn
	//当前执行的命令
	cmds []string
	//当前命令有几节
	cmdLength int
	//会话结束标志 true结束、false未结束
	endFlag bool
}

/**************************验证并解析身份**********************************/
func (ctxt *context) addFolder() {

}

/**************************验证并解析身份**********************************/
func (ctxt *context) addUser() { //2009权限不足、2007用户信息重复、2002命令格式不对、4500服务器出错、4005添加数据出错、1001成功
	var respCode string
	var err error
	defer func() {
		respByte(ctxt, respCode)
	}()
	//判断是不是管理员。
	if !ctxt.isAdmin() {
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
func (ctxt *context) login() { //身份验证成功返回1002；程序出错返回4500；登录命令错误返回2006；数据库查询出错2008;
	var rst string
	var serverErr string
	defer func() {
		if serverErr != "" {
			respByte(ctxt, serverErr)
		} else if rst != "" {
			respByte(ctxt, rst)
			loging.Loger.Info("[%s:登陆成功]\n", ctxt.sessionInfo.UserIp)
		}
	}()
	//--------------------------------------------------
	//1、检查命令格式是否正确。不正确：2006
	if len(ctxt.cmds) != 3 {
		serverErr = recode.RECODE_SESSIONERR
		return
	}
	//2、检查命令是否正确。不正确：2006
	if ctxt.cmds[0] == "LOGIN" {
		//2-1查询用户是否存在。不存在：2008  出错4001
		ctxt.userInfo.Name = ctxt.cmds[1]
		err := dbops.DBContext.Read(&ctxt.userInfo, "Name")
		if err != nil {
			if err == orm.ErrNoRows {
				loging.Loger.Error("该用户不存在", err)
				serverErr = recode.RECODE_USERERR
				return
			}
			serverErr = recode.RECODE_DBERR
			return
		}
		//2-2校验密码是否正确，如果正确生成一个sessionID。  出错：4500，登录成功：1002
		pwd, err := com.GetMD5(ctxt.cmds[2])
		if err != nil {
			loging.Loger.Error("MD5生成出错", err)
			serverErr = recode.RECODE_SERVERERR
			return
		}
		if pwd != ctxt.userInfo.PassWord {
			serverErr = recode.RECODE_PWDERR
			return
		}
		uuid, err := com.GetUUID()
		if err != nil {
			loging.Loger.Error("uuid生成出错", err)
			serverErr = recode.RECODE_SERVERERR
			return
		}
		ctxt.sessionInfo.SId = uuid
		rst = recode.RECODE_LOGINOK

		//2-3登录成功后 加载用户home目录信息。
		dbops.DBContext.QueryTable("FileInfo").Filter("name__contains", ctxt.userInfo.Name).
			Filter("user__id", ctxt.userInfo.Id).
			One(&ctxt.currentFolder)
		if ctxt.currentFolder.Id < 1 {
			loging.Loger.Error("加载用户%s 的home目录信息失败", ctxt.sessionInfo.UserIp)
		}
	} else {
		serverErr = recode.RECODE_SESSIONERR
		return
	}

}

/**************************未识别的命令 **********************************/
func (ctxt *context) unIdentified() {
	ctxt.conn.Write([]byte(recode.RECODE_UNIDENTIFIED))
}

/**************************拒绝连接***************************************/
func (ctxt *context) turnDown() {
	ctxt.endFlag = true
	ctxt.conn.Write([]byte(recode.RECODE_PARAMERR))

}

/**************************退出指令***************************************/
func (ctxt *context) exit() {
	ctxt.endFlag = true
	ctxt.conn.Write([]byte(recode.RECODE_OK))
}

/*************************************************************************/
/*************************************************************************/
/*************************************************************************/

//辅助方法：判断是否为管理员账户
func (ctxt *context) isAdmin() bool {
	for i, _ := range adminList {
		if adminList[i] == ctxt.userInfo.Name {
			return true
		}
	}
	return false
}
