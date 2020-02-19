package httpreq

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"mzy-spider/model"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultJS     = "assets/js/custom/config.js"
	DownLoadUrls  = []string{"亚洲电影", "欧美电影", "制服丝袜", "强奸乱伦", "国产自拍", "变态另类", "经典三级", "成人动漫"} // 下载专区 菜单
	VideoMenuList = []string{"亚洲无码", "女优专辑", "短视频", "国产精品", "中文字幕", "欧美精品", "成人动漫"}          // 在线视频 菜单
	NewURL        = ""
	defaultHome   = "index/home.html"
	IsNext        = false
)

func init() {

}

// 猫咪视频资源(所有)
func Run(addr string) {
	if addr == "" {
		panic(errors.New("string empty"))
	}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"),
	)
	c.SetRequestTimeout(5 * time.Second)
	// 生成新的操作对象
	videoColly, videoListColly, pageSizeColly, videoDetailColly := c.Clone(), c.Clone(), c.Clone(), c.Clone()

	c.OnHTML("script", func(element *colly.HTMLElement) {
		var a = element.Attr("src")
		if strings.Contains(a, defaultJS) {
			c.Visit(addr + a) // 获取最新域名
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(response *colly.Response, e error) {
		until.PrintlnMsg(true, true, fmt.Sprintf("ERROR CODE : %d,%s,%s", response.StatusCode, response.Request.URL.String(), e.Error()))
	})
	c.OnResponse(func(response *colly.Response) {
		if strings.Contains(response.Ctx.Get("url"), "assets") {
			var (
				body       = string(response.Body)
				firstIndex = strings.Index(body, "window.default_line = \"")
				lastIndex  = strings.LastIndex(body, "\";")
				newURL     = strings.ReplaceAll(body[firstIndex:lastIndex], "window.default_line = \"", "")
			)
			videoColly.Visit(newURL + defaultHome)
		}
	})

	videoColly.OnRequest(func(request *colly.Request) {
		url := request.URL.String()
		NewURL = request.URL.Host
		fmt.Println("videoColly's Visiting", url)
	})

	videoColly.OnResponse(func(response *colly.Response) {
		var url = response.Request.URL.String() + "/" + defaultHome
		if !strings.Contains(url, NewURL) {
			videoColly.Visit(url)
		}
	})
	// 获取菜单链接
	videoColly.OnHTML("ul.row-item-content > li > a", func(htmlElement *colly.HTMLElement) {
		var (
			url  = htmlElement.Attr("href")
			text = strings.ReplaceAll(htmlElement.Text, " ", "")
			host = htmlElement.Request.URL.Scheme + "://" + htmlElement.Request.URL.Host
		)
		// 单独获取在线资源
		if !IsNext {
			for _, v := range VideoMenuList {
				t := fmt.Sprintf("/shipin/list-%s.html", v)
				pageSizeColly.Visit(host + t)
			}
		}
		IsNext = true

		if containsKey(text) {
			pageSizeColly.Visit(host + url)
		}
	})
	// 获取分页列表
	pageSizeColly.OnHTML(".pagination > a:nth-last-child(4)", func(htmlElement *colly.HTMLElement) {
		var pageSize, _ = strconv.ParseInt(htmlElement.Text, 0, 64)
		for i := 1; i <= int(pageSize); i++ {
			path := htmlElement.Request.URL.Path[0:strings.Index(htmlElement.Request.URL.Path, ".html")]
			path = path + fmt.Sprintf("-%d.html", i)
			url := htmlElement.Request.URL.Scheme + "://" + htmlElement.Request.URL.Host + path
			t := time.Now()
			err := videoListColly.Visit(url)
			t1 := time.Now()
			fmt.Println(url + fmt.Sprintf(" -- %v", t1.Sub(t)))
			if err != nil {
				continue
			}
		}
	})
	// 获取链接
	videoListColly.OnHTML("#tpl-img-content > li > a", func(htmlElement *colly.HTMLElement) {
		var url = htmlElement.Attr("href")
		title := htmlElement.Attr("title")
		if IsExist(title) {
			videoDetailColly.Visit(htmlElement.Request.URL.Scheme + "://" + htmlElement.Request.URL.Host + url)
		}
	})
	// 短视频专区
	videoListColly.OnHTML("#grid > li > a", func(htmlElement *colly.HTMLElement) {
		var url = htmlElement.Attr("href")
		title := htmlElement.Attr("title")

		if IsExist(title) {
			videoDetailColly.Visit(htmlElement.Request.URL.Scheme + "://" + htmlElement.Request.URL.Host + url)
		}
	})

	// 女优专区链接
	videoDetailColly.OnHTML("#tpl-img-content > li > a", func(element *colly.HTMLElement) {
		var url = element.Attr("href")
		title := element.Attr("title")
		if IsExist(title) {
			videoDetailColly.Visit(element.Request.URL.Scheme + "://" + element.Request.URL.Host + url)
		}
	})

	// 获取视频详情
	videoDetailColly.OnHTML("#main-container", func(htmlElement *colly.HTMLElement) {
		if selection := htmlElement.DOM.Find("#tpl-img-content"); len(selection.Nodes) > 0 { // 女优专区
		} else {
			down0Query := htmlElement.DOM.Find("#lin1k0")
			down1Query := htmlElement.DOM.Find("#lin1k1")
			titleQuery := htmlElement.DOM.Find(".row > h2")
			var (
				mysqlDB           = stock.ActionMysql.Db
				video0DownLoad, _ = down0Query.Attr("value")
				video1DownLoad, _ = down1Query.Attr("value")
				title             = titleQuery.Text()
				htmlURL           = htmlElement.Request.URL.String()
			)
			finish := 1
			if video0DownLoad == "" || video1DownLoad == "" || title == "" || htmlURL == "" {
				finish = 0
			}
			ret, err := mysqlDB.Exec("INSERT movie_info SET title=?,html_url=?,dateline=?,video_url=?,thunder_url=?,finish=?,is_down=0", title, htmlURL, time.Now().Unix(), video0DownLoad, video1DownLoad, finish)
			if err != nil {
				panic(err.Error())
			}
			if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
				log.Println("[写入数据库失败]:", ok.Error())
			}
			fmt.Println("Write -> " + title)
		}
	})

	c.Visit(addr)
}
func IsExist(title string) bool {
	mysqlDB := stock.ActionMysql.Db
	exist := int64(0)
	err := mysqlDB.QueryRow("SELECT COUNT(1) FROM movie_info WHERE title=? ", title).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return exist == 0
}
func containsKey(str string) bool {
	if str != "" {
		for _, v := range DownLoadUrls {
			if strings.Contains(v, str) {
				return true
			}
		}
	}
	return false
}

// 批量下载迅雷X
func DownloadVideo(resp http.ResponseWriter, PageSize int64) {
	if PageSize == 0 {
		PageSize = 5
	}
	var (
		mysqlDB = stock.ActionMysql.Db
		task    = model.Thunder{
			MinVersion:         "10.0.1.0",
			TaskGroupName:      "视频资源",
			ThreadCount:        10,
			ThunderInstallPack: "http://down.sandai.net/thunderx/XunLeiSetup10.1.1.148Beta.exe",
		}
		ids []string
	)
	rows, err := mysqlDB.Query("SELECT id,thunder_url,title,video_url FROM movie_info WHERE is_down=0 LIMIT ?", PageSize)
	if err != nil && err != sql.ErrNoRows {
		panic(err.Error())
	}
	for rows.Next() {
		id := int64(0)
		url := ""
		m := model.Task{}
		err = rows.Scan(&id, &m.Url, &m.Name, &url)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if strings.Contains(url, "rmvb") {
			m.Name = m.Name + ".rmvb"
		} else {
			m.Name = m.Name + ".mp4"
		}
		ids = append(ids, strconv.Itoa(int(id)))
		task.Tasks = append(task.Tasks, m)
	}
	rows.Close()

	if len(ids) > 0 {
		str := fmt.Sprintf("UPDATE movie_info SET is_down=1 WHERE id IN(%s)", strings.Join(ids, ","))
		ret, err := mysqlDB.Exec(str)
		if err != nil {
			panic(err.Error())
		}
		if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
			log.Println("[写入数据库失败]:", ok.Error())
		}
	}

	until.Json(resp, task)
}
