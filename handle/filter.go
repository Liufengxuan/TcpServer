package handle

import (
	"TcpServer/models/dbops"
	"TcpServer/models/recode"
	"fmt"
	"log"
	"net"
	"strings"
)

var lineBreakLength int

type currentInfo struct {
	userInfo    dbops.User
	sessionInfo dbops.Session
}

func HandlerConn(conn net.Conn) {
	var curInfo currentInfo
	defer conn.Close()
	//读消息。
	log.Println("处理", conn.RemoteAddr().String())
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("[读取用户内容时出现异常：%s]\n", err)
		return
	}

	//处理分割消息。
	cmd := string(buf[:n-lineBreakLength])

	fmt.Println(lineBreakLength)

	cmds := strings.Split(cmd, " ")
	if len(cmds) > 1 {
		if auth(cmds, &curInfo) {
			log.Printf("[ok(%s)]\n", cmds)
		}

	} else {
		conn.Write([]byte(recode.RECODE_PARAMERR))
	}

}
