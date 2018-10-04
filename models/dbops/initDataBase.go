package dbops

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" // import your used driver
)

var DBContext orm.Ormer

func init() {
	// set default database
	orm.RegisterDataBase("default", "sqlite3", "./Data.db", 30)

	// register model
	orm.RegisterModel(new(User), new(Session), new(FileInfo))

	// create table
	orm.RunSyncdb("default", false, false)
	DBContext = orm.NewOrm()
}
