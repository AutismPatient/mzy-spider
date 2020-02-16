package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"mzy-spider/httpreq"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"time"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		var (
			mysqlDB = stock.ActionMysql.Db
			runKey  = req.URL.Query().Get("run_key")
		)
		rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
		if err != nil || runKey != rely {
			resp.Write([]byte("INVALID PARAMETER VALUE"))
			return
		}
		ret, err := mysqlDB.Exec("INSERT runtime SET run_key=?,is_run=?,run_time=?", runKey, 1, time.Now().Unix())
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}

		httpreq.Run("https://www.maomiav.com/")

		ret, err = mysqlDB.Exec("UPDATE runtime SET is_run=? WHERE run_key=?", 0, runKey)
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}
		log.Println("任务执行完毕")
		resp.Write([]byte("end of run"))
	}
	resp.WriteHeader(404)
}

func main() {

	addr := ":8888"

	mux := http.NewServeMux()
	mux.HandleFunc("/program/run", handler)

	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    5 * time.Second,
	}

	err := srv.ListenAndServe()

	if err != nil {
		panic(err)
	}

	_, err = stock.Redis.Do("set", "run_key", until.GenerateToken(0))
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Now().Format("2006-05-04 15:30:03"), "站点初始化成功")
}
