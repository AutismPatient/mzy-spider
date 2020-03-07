package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mzy-spider/httpreq"
	"mzy-spider/reg"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/url"
)

type VideoDownloadController struct {
}

func init() {
	reg.Register("download", &VideoDownloadController{}, true)
}

func (c *VideoDownloadController) Default(r *gin.Context) {
	c.Run(r)
}

func (c *VideoDownloadController) Run(r *gin.Context) {
	var (
		menu, _   = url.PathUnescape(r.Query("menu"))
		search, _ = url.PathUnescape(r.Query("q"))
		size      = r.GetInt64("page_size")
	)
	httpreq.DownloadVideo(r.Writer, size, menu, search)
}

func (c *VideoDownloadController) Search(r *gin.Context) {
	var (
		action    = stock.ActionMysql.Db
		search, _ = url.PathUnescape(r.Query("q"))
		page      = r.GetInt64("page")
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
	until.Json(r.Writer, list)
}

func (c *VideoDownloadController) SearchRun(r *gin.Context) {
	var (
		movies = r.Query("ids")
	)
	httpreq.DownloadVideoByIDS(r.Writer, movies)
}
