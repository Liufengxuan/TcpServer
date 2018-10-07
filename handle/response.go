package handle

import (
	"TcpServer/models/recode"
)

var sysError = "4"
var usrError = "2"
var ok = "1"

func respByte(curInfo *context, code string) {
	//log.Println(code)
	State := code[0:1]
	resp := []byte(code)
	switch State {
	case sysError:
		curInfo.conn.Write(resp)
		curInfo.endFlag = true
	case usrError:
		curInfo.conn.Write(resp)
	case ok:
		curInfo.conn.Write(resp)
	default:
		curInfo.conn.Write([]byte(recode.RECODE_UNKNOWERR))
	}

	return
}
