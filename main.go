package main

import (
	"database/sql"
	"fmt"
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

const IsUPDATEREADY = false //是否更新密钥,仅用于调试模式

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

		defer func() {
			httpreq.Run("https://www.maomiav.com/")
			fmt.Println(time.Now().Format("2006-01-02 15:04:05") + "-- 任务执行完成.")
		}()

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
		runKey    = req.URL.Query().Get("run_key")
		menu, _   = url.PathUnescape(req.URL.Query().Get("menu"))
		search, _ = url.PathUnescape(req.URL.Query().Get("q"))
		size, _   = strconv.ParseInt(req.URL.Query().Get("page_size"), 0, 64)
	)
	rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
	if err != nil || runKey != rely {
		resp.Write([]byte("INVALID PARAMETER VALUE"))
		return
	}
	httpreq.DownloadVideo(resp, size, menu, search)
}
func workDownLoadByIDSHandle(resp http.ResponseWriter, req *http.Request) {
	var (
		runKey = req.URL.Query().Get("run_key")
		movies = req.URL.Query().Get("ids")
	)
	rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
	if err != nil || runKey != rely {
		resp.Write([]byte("INVALID PARAMETER VALUE"))
		return
	}
	httpreq.DownloadVideoByIDS(resp, movies)
}
func searchMovieHandle(resp http.ResponseWriter, req *http.Request) {
	var (
		runKey    = req.URL.Query().Get("run_key")
		action    = stock.ActionMysql.Db
		search, _ = url.PathUnescape(req.URL.Query().Get("q"))
		page, _   = strconv.ParseInt(req.URL.Query().Get("page"), 0, 64)
		where     = ""
		offset    = int64(0)
	)
	type MovieSub struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		DownLoadUrl string `json:"download_url"`
		HtmlPath    string `json:"html_path"`
		Menu        string `json:"menu"`
	}
	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * 15
	rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
	if err != nil || runKey != rely {
		resp.Write([]byte("INVALID PARAMETER VALUE"))
		return
	}
	if search != "" {
		where = where + fmt.Sprintf(" AND MATCH(title,menu) AGAINST('*%s*'IN BOOLEAN MODE)", search)
	}
	rows, err := action.Query("SELECT id,thunder_url,title,html_url,menu FROM movie_info WHERE is_down=0"+where+" ORDER BY dateline DESC LIMIT ?,?", offset, 15)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	var list = []MovieSub{}
	for rows.Next() {
		l := MovieSub{}
		err = rows.Scan(&l.Id, &l.DownLoadUrl, &l.Title, &l.HtmlPath, &l.Menu)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		list = append(list, l)
	}
	rows.Close()
	until.Json(resp, list)
}
func htmlDownLoadHandle(resp http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./html/download.html")
	if err != nil {
		log.Println("err:", err)
		return
	}
	t.Execute(resp, nil)
}
func htmlSearchHandle(resp http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./html/search.html")
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
	mux.HandleFunc("/download/search", searchMovieHandle)
	mux.HandleFunc("/download/searchHtml", htmlSearchHandle)
	mux.HandleFunc("/download/search_run", workDownLoadByIDSHandle)

	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    5 * time.Second,
	}

	if IsUPDATEREADY {
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
	// 判断是否启用快捷模式
	actionMysql := stock.ActionMysql.Db
	dateLine := int64(0)

	err := actionMysql.QueryRow("SELECT dateline FROM movie_info ORDER BY dateline DESC LIMIT 1").Scan(&dateLine)
	if err != nil {
		panic(err)
	}

	if (time.Now().Unix() - dateLine) < (86400 * 5) { // 5天内启用
		_, err := stock.Redis.Do("set", "next", time.Now().Unix())
		if err != nil {
			panic(err)
		}
		log.Println(time.Now().Format("2006-01-02 15:04:05") + "快捷模式")
	}

	until.PrintlnMsg(false, true, time.Now().Format("2006-01-02 15:04:05")+" 站点初始化成功，秘钥已更新")

	log.Fatal(srv.ListenAndServe())
}
