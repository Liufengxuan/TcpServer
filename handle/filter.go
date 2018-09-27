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
		log.Println("——>处理结束", userContext.conn.RemoteAddr().String())
		userContext.conn.Close()

	}()

	log.Println("——>开始处理", userContext.conn.RemoteAddr().String())
	rBuf := make([]byte, readBufferSize)

	for {
		if userContext.endFlag {
			return
		}

		n, err := userContext.conn.Read(rBuf)
		if err != nil {
			log.Printf("[读取用户内容时出现异常,可能由于用户断开了连接：%s]\n", err)
			return
		}
		//处理分割消息。
		cmd := string(rBuf[:n-lineBreakLength])
		userContext.cmds = strings.Split(cmd, " ")

		//-------------------------------------------------------------
		if len(userContext.cmds) > 1 {
			cmdFilter(&userContext)
		} else {
			cmdFilter(&userContext)
		}

	}

}

/**************************method*********************************************************/

func cmdFilter(userContext *context) {
	cmd := strings.ToUpper(userContext.cmds[0])
	switch cmd {
	case "LOGIN":
		auth(userContext)
	case "EXIT":
		exit(userContext)
	default:
		unIdentified(userContext)
	}
}
