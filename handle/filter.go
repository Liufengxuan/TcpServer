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
	defer func() {
		loging.Loger.Info("[%s:会话结束]", userContext.conn.RemoteAddr().String())
		userContext.conn.Close()
	}()

	loging.Loger.Info("[%s:已连接]", userContext.conn.RemoteAddr().String())

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
			cmd := com.CharReplace(msg)
			userContext.cmds = strings.Split(cmd, " ")
			userContext.cmdLength = len(userContext.cmds)
			cmdRoute(&userContext)
		case err := <-errCh:
			loging.Loger.Warning("[读取用户消息时出现异常：%s]\n", err)
			return
		case <-time.After(respTimeout):
			loging.Loger.Info("[%s:用户长时间未响应]", userContext.conn.RemoteAddr().String())
			return
		}

	}
	//---------------------------------------------------------

}

/**************************method*********************************************************/
func cmdRoute(userContext *context) {
	userContext.cmds[0] = strings.ToUpper(userContext.cmds[0])

	//判断用户是否为第一次登陆
	if userContext.sessionInfo.SId == "" {
		userContext.auth()
	} else {
		//————————————————————————————————————————————————
		switch userContext.cmds[0] {
		//EXIT 命令
		case "QUIT":
			userContext.exit()
		//CREATE命令
		case "CREATE":
			if len(userContext.cmds) >= 3 {
				userContext.cmds[1] = strings.ToUpper(userContext.cmds[1])
				switch userContext.cmds[1] {
				//USER命令
				case "USER":
					userContext.addUser()
					//DIR命令
				//case "DIR":
				default:
					userContext.unIdentified()
				}
			} else {
				userContext.unIdentified()
			}
		//没有的命令
		default:
			userContext.unIdentified()
		}
		//————————————————————————————————————————————————
	}

}

func isAdmin(s string) bool {
	for i, _ := range adminList {
		if adminList[i] == s {
			return true
		}
	}
	return false
}
