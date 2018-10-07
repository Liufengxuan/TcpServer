package handle

import (
	"TcpServer/com"
	"TcpServer/loging"
	"net"
	"strings"
	"time"
)

/**************************method*********************************************************/

//HandlerConn 调用此方法处理连接。
func HandlerConn(conn net.Conn) {
	var userContext context
	userContext.conn = conn
	userContext.endFlag = false
	userContext.sessionInfo.UserIp = userContext.conn.RemoteAddr().String()
	defer func() {
		loging.Loger.Info("[%s:会话结束]", userContext.sessionInfo.UserIp)
		userContext.conn.Close()
	}()

	loging.Loger.Info("[%s:已连接]", userContext.sessionInfo.UserIp)

	//--------------------------------------------------------------
	for {
		//通过结束标记 判断此会话是否已经结束
		if userContext.endFlag {
			return
		}

		//接收和 解析消息
		msgCh := make(chan string)
		errCh := make(chan error)
		go func() {
			rBuf := make([]byte, readBufferSize)
			n, err := userContext.conn.Read(rBuf)
			if err != nil {
				errCh <- err
			}
			msgCh <- string(rBuf[:n])
		}()

		select {
		case msg := <-msgCh:
			userContext.cmds = com.CmdFormat(msg)
			userContext.cmdLength = len(userContext.cmds)
			cmdRoute(&userContext)
		case err := <-errCh:
			loging.Loger.Warning("[读取用户消息时出现异常：%s]\n", err)
			return
		case <-time.After(respTimeout):
			loging.Loger.Info("[%s:用户长时间未响应]", userContext.sessionInfo.UserIp)
			return
		}

	}
	//---------------------------------------------------------

}

/**************************method*********************************************************/
func cmdRoute(userContext *context) {
	if userContext.cmdLength > 0 {
		userContext.cmds[0] = strings.ToUpper(userContext.cmds[0])
	} else {
		userContext.unIdentified()
		return
	}

	//判断用户是否登录
	if userContext.sessionInfo.SId == "" {
		userContext.login()
	} else {
		//————————————————————————————————————————————————
		switch userContext.cmds[0] {
		//EXIT 命令
		case "QUIT":
			userContext.exit()
		//CREATE命令
		case "CREATE":
			createCmd(userContext)
		//没有的命令
		default:
			userContext.unIdentified()
		}
		//————————————————————————————————————————————————
	}

}

//**************************当首级指令为create执行此方法***************************************************
func createCmd(userContext *context) {
	if len(userContext.cmds) >= 3 {
		userContext.cmds[1] = strings.ToUpper(userContext.cmds[1])
		switch userContext.cmds[1] {
		//USER命令
		case "USER":
			userContext.addUser()
			//DIR命令
		case "DIR":
			userContext.addFolder()
		default:
			userContext.unIdentified()
		}
	} else {
		userContext.unIdentified()
	}
}
