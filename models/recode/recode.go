package recode

//4000 开头系统级错误、断开用户的连接
//2000 开头用户方错误、不断开连接
//1000 开头正确状态。

const (
	RECODE_OK      = "1001"
	RECODE_LOGINOK = "1002"

	RECODE_UNIDENTIFIED  = "2001"
	RECODE_CMDFORMATERR  = "2002"
	RECODE_PWDERR        = "2003"
	RECODE_LOGINPARAMERR = "2004"
	RECODE_LOGINERR      = "2005"
	RECODE_SESSIONERR    = "2006"
	RECODE_REPEATUSER    = "2007"
	RECODE_USERERR       = "2008"

	RECODE_DBERR      = "4001"
	RECODE_ADDDATAERR = "4005"
	RECODE_NODATA     = "4002"
	RECODE_DATAEXIST  = "4003"
	RECODE_DATAERR    = "4004"
	RECODE_PARAMERR   = "4103"
	RECODE_ROLEERR    = "4105"
	RECODE_REQERR     = "4201"
	RECODE_IPERR      = "4202"
	RECODE_THIRDERR   = "4301"
	RECODE_IOERR      = "4302"
	RECODE_SERVERERR  = "4500"
	RECODE_UNKNOWERR  = "5000" //"未知错误",
)

var recodeText = map[string]string{
	RECODE_OK:        "成功",
	RECODE_DBERR:     "数据库查询错误",
	RECODE_NODATA:    "无数据",
	RECODE_DATAEXIST: "数据已存在",
	RECODE_DATAERR:   "数据错误",
	//------------------------------------------
	RECODE_LOGINOK:    "登陆成功",
	RECODE_SESSIONERR: "用户未登录",
	RECODE_LOGINERR:   "用户登录失败",
	RECODE_PARAMERR:   "参数错误",
	RECODE_USERERR:    "用户不存在或未激活",
	RECODE_ROLEERR:    "用户身份错误",
	RECODE_PWDERR:     "密码错误",
	RECODE_REPEATUSER: "用户已存在",

	//------------------------------------------------
	RECODE_REQERR:     "非法请求或请求次数受限",
	RECODE_IPERR:      "IP受限",
	RECODE_THIRDERR:   "第三方系统错误",
	RECODE_IOERR:      "文件读写错误",
	RECODE_SERVERERR:  "内部错误",
	RECODE_UNKNOWERR:  "未知错误",
	RECODE_ADDDATAERR: "添加数据出错",
	//---------------------------------------------------------
	RECODE_CMDFORMATERR:  "命令格式错误",
	RECODE_UNIDENTIFIED:  "未识别出该命令",
	RECODE_LOGINPARAMERR: "登陆参数错误",
}

func RecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKNOWERR]
}
