package dbops

//User 用户表
type User struct {
	Id       int        `json:"id"`
	Name     string     `orm:"size(32);unique"  json:"name"`
	PassWord string     `orm:"size(64)" json:"password"`
	Session  []*Session `orm:"reverse(many)"`
}

//Session 会话列表。
type Session struct {
	Id         int
	SId        string `orm:"size(64)"`
	BreakPoint int    `orm:"default(0)"`
	FileId     int
	UserIp     string `orm:"size(64)"`
	User       *User  `orm:"rel(fk)"`
}
