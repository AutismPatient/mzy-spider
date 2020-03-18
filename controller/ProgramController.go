package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mzy-spider/httpreq"
	"mzy-spider/reg"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"time"
)

type ProgramController struct {
}

func init() {
	reg.Register("program", &ProgramController{}, true)
}
func (c *ProgramController) Default(r *gin.Context) {
	c.Run(r)
}
func (c *ProgramController) Run(r *gin.Context) {
	if r.Request.Method == http.MethodGet {
		var (
			mysqlDB = stock.ActionMysql.Db
			runKey  = r.Query("run_key")
		)
		ret, err := mysqlDB.Exec("INSERT runtime SET run_key=?,is_run=?,run_time=?", runKey, 1, time.Now().Unix())
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}

		defer func() {
			httpreq.Run(httpreq.ConfigValue["地址"])
			fmt.Println(until.TimeFormat(time.Now()) + "-- 任务执行完成.")
		}()

		ret, err = mysqlDB.Exec("UPDATE runtime SET is_run=? WHERE run_key=?", 0, runKey)
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}
		r.Writer.Write([]byte("Task begins"))
		return
	}
	r.Writer.WriteHeader(http.StatusNotFound)
}
