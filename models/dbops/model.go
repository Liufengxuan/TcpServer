package dbops

//User 用户表
type User struct {
	Id       int         `json:"id"`
	Name     string      `orm:"size(32);unique"  json:"name"`
	PassWord string      `orm:"size(64)" json:"password"`
	Session  []*Session  `orm:"reverse(many);null"`
	File     []*FileInfo `orm:"reverse(many);null"`
}

//Session 会话列表。
type Session struct {
	Id         int
	SId        string `orm:"size(64)"`
	BreakPoint int    `orm:"default(0);null"`
	FileId     int
	UserIp     string `orm:"size(64);null"`
	User       *User  `orm:"rel(fk);null"`
}
type FileInfo struct {
	Id       int
	User     *User  `orm:"rel(fk)"`
	Size     int    `orm:"null"`
	IsDir    bool   `orm:"default(false)"`
	IsHome   bool   `orm:"default(false)"`
	Name     string `orm:"size(64)"`
	ParentId int    `orm:"default(0);null"`
	Md5      string `orm:"size(64)"`
}
