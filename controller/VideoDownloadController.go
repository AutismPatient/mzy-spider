package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mzy-spider/httpreq"
	"mzy-spider/reg"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/url"
	"strconv"
	"strings"
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

// GetQueryOfInt64 获取request参数 INT64
func GetQueryOfInt64(r *gin.Context, key string) (it int64) {
	q := r.Query(key)
	it, err := strconv.ParseInt(q, 0, 64)
	if err != nil {
		return 0
	}
	return it
}
func (c *VideoDownloadController) Search(r *gin.Context) {
	var (
		action    = stock.ActionMysql.Db
		search, _ = url.PathUnescape(r.Query("q"))
		page      = GetQueryOfInt64(r, "page")
		isDown    = GetQueryOfInt64(r, "is_down")
		where     = ""
		offset    = int64(0)
	)
	type MovieSub struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		DownLoadUrl string `json:"download_url"`
		HtmlPath    string `json:"html_path"`
		Menu        string `json:"menu"`
		IsDown      int64  `json:"is_down"`
	}
	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * 15
	if search != "" {
		where = where + fmt.Sprintf(" AND MATCH(title,menu) AGAINST('*%s*'IN BOOLEAN MODE)", search)
	}
	var orderBy = "dateline"
	if isDown == 1 {
		orderBy = "down_date"
	}
	var str = fmt.Sprintf("SELECT id,thunder_url,title,html_url,menu FROM movie_info WHERE is_down=%d%s ORDER BY %s DESC LIMIT %d,%d", isDown, where, orderBy, offset, 15)
	fmt.Println(str)
	rows, err := action.Query(str)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	var list []MovieSub
	for rows.Next() {
		l := MovieSub{}
		err = rows.Scan(&l.Id, &l.DownLoadUrl, &l.Title, &l.HtmlPath, &l.Menu)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		l.IsDown = isDown
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

// PlayVideo 视频播放
func (c *VideoDownloadController) PlayVideo(r *gin.Context) {
	var (
		Name = r.Query("title")
		Dirs = []string{until.C, until.D, until.E, until.F}
		Path = ""
	)
	Name = strings.Replace(Name, " ", "", -1)
	for _, v := range Dirs {
		if Path == "" {
			FileInfos, _ := ioutil.ReadDir(v)
			for _, v1 := range FileInfos {
				fn := v1.Name()[:strings.LastIndex(v1.Name(), ".")]
				if Name == fn {
					Path = "/" + strings.Replace(v, ":/视频资源", "", -1) + "/" + v1.Name()
					break
				}
			}
		}
	}
	r.String(200, Path)
}
