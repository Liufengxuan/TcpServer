package handle

import "TcpServer/models/recode"

func auth(curInfo *context) {
	curInfo.conn.Write([]byte("成功" + curInfo.cmds[1] + curInfo.cmds[2]))
}

func unIdentified(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_UNIDENTIFIED))
}

/************************************************************/

func turnDown(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_PARAMERR))
	curInfo.endFlag = true
}

func exit(curInfo *context) {
	curInfo.conn.Write([]byte(recode.RECODE_OK))
	curInfo.endFlag = true
}
