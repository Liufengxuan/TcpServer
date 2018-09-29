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
	conn        net.Conn
	cmds        []string
	endFlag     bool
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
		if len(userContext.cmds) > 1 {
			cmdFilter(&userContext)
		} else {
			cmdFilter(&userContext)
		}
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
		case "EXIT":
			exit(userContext)
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
