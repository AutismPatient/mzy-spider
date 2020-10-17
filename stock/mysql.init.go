package stock

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Mysql struct {
	Db *sql.DB
}

var (
	ActionMysql = &Mysql{}
)

func init() {
	ActionMysql.newMysqlConn()
}
func (m *Mysql) newMysqlConn() {
	var source = fmt.Sprintf("root:qazwsxedcR178@tcp(127.0.0.1:3306)/movie")
	db, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}
	if e := db.Ping(); e != nil {
		panic(e.Error())
	}
	db.SetMaxOpenConns(120)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(20 * time.Second)
	m.Db = db
}
