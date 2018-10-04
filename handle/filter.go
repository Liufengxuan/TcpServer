package handle

import (
	"TcpServer/models/dbops"
	"log"
	"net"
	"strings"
)

type context struct {
	userInfo    dbops.User
	sessionInfo dbops.Session
	//当前用户的   连接对象
	conn net.Conn
	//当前用户的    命令
	cmds []string
	//当前用户的    会话结束标志
	endFlag bool
}

/**************************method*********************************************************/

//HandlerConn 调用此方法处理连接。
func HandlerConn(conn net.Conn) {
	var userContext context
	userContext.conn = conn
	userContext.endFlag = false
	defer func() {
		log.Printf("[%s:已断开]", userContext.conn.RemoteAddr().String())
		userContext.conn.Close()
	}()

	log.Printf("[%s:已连接]", userContext.conn.RemoteAddr().String())
	rBuf := make([]byte, readBufferSize)

	//--------------------------------------------------------------
	for {
		//通过结束标记 判断此会话是否已经结束
		if userContext.endFlag {
			return
		}

		//接收和 解析消息
		n, err := userContext.conn.Read(rBuf)
		if err != nil {
			log.Printf("[读取用户内容时出现异常,可能由于用户断开了连接：%s]\n", err)
			return
		}
		cmd := charReplace(string(rBuf[:n]))
		userContext.cmds = strings.Split(cmd, " ")

		//通过过滤命令来处理消息。
		cmdFilter(&userContext)
	}
	//---------------------------------------------------------

}

/**************************method*********************************************************/

func cmdFilter(userContext *context) {
	userContext.cmds[0] = strings.ToUpper(userContext.cmds[0])

	//判断用户是否为第一次登陆
	if userContext.sessionInfo.SId == "" {
		auth(userContext)
	} else {
		//————————————————————————————————————————————————
		switch userContext.cmds[0] {
		//EXIT 命令
		case "EXIT":
			exit(userContext)
		//CREATE命令
		case "CREATE":
			if len(userContext.cmds) >= 3 {
				userContext.cmds[1] = strings.ToUpper(userContext.cmds[1])
				switch userContext.cmds[1] {
				//创建用户
				case "USER":
					addUser(userContext)
				case "DIR":
				default:
					unIdentified(userContext)
				}
			} else {
				unIdentified(userContext)
			}
		//没有的命令
		default:
			unIdentified(userContext)
		}
		//————————————————————————————————————————————————
	}

}

/**************************method*********************************************************/

func charReplace(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}
