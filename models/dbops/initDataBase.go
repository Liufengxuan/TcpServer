package dbops

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" // import your used driver
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "sqlite3", "./$database/Data.db", 30)

	// register model
	orm.RegisterModel(new(User), new(Session))

	// create table
	orm.RunSyncdb("default", false, false)
}
