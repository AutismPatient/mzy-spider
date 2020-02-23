package main

import (
	"github.com/gomodule/redigo/redis"
	"html/template"
	"log"
	"mzy-spider/httpreq"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const IsUPDATEKEY = true //是否更新密钥,仅用于调试模式

func workHandle(resp http.ResponseWriter, req *http.Request) {
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

		defer httpreq.Run("https://www.maomiav.com/")

		ret, err = mysqlDB.Exec("UPDATE runtime SET is_run=? WHERE run_key=?", 0, runKey)
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}
		resp.Write([]byte("Task begins"))
		return
	}
	resp.WriteHeader(404)
}
func workDownLoadHandle(resp http.ResponseWriter, req *http.Request) {
	var (
		runKey  = req.URL.Query().Get("run_key")
		menu, _ = url.PathUnescape(req.URL.Query().Get("menu"))
		size, _ = strconv.ParseInt(req.URL.Query().Get("page_size"), 0, 64)
	)
	rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
	if err != nil || runKey != rely {
		resp.Write([]byte("INVALID PARAMETER VALUE"))
		return
	}
	httpreq.DownloadVideo(resp, size, menu)
}
func htmlDownLoadHandle(resp http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./html/download.html")
	if err != nil {
		log.Println("err:", err)
		return
	}
	t.Execute(resp, nil)
}
func main() {

	addr := ":8888"

	mux := http.NewServeMux()

	mux.HandleFunc("/program/run", workHandle)
	mux.HandleFunc("/download/run", workDownLoadHandle)
	mux.HandleFunc("/download/index", htmlDownLoadHandle)

	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    5 * time.Second,
	}

	if IsUPDATEKEY {
		token := until.GenerateToken(0)
		_, err := stock.Redis.Do("set", "run_key", token)
		if err != nil {
			panic(err)
		}

		err = until.SendEmail([]string{"1010014622@qq.com", "1746793113@qq.com"}, "run_key", token)

		if err != nil {
			panic(err)
		}
	}

	until.PrintlnMsg(false, true, time.Now().Format("2006-01-02 15:04:05")+" 站点初始化成功，秘钥已更新")

	log.Fatal(srv.ListenAndServe())
}
