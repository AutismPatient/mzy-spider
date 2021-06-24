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
	//PassWord = "!Mngsy4f6O,?0sj"
	PassWord = "123456"
)

func init() {
	ActionMysql.newMysqlConn()
}
func (m *Mysql) newMysqlConn() {
	var source = fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/movie", PassWord)
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
