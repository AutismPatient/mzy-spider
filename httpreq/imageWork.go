package httpreq

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	base   = "http://www.mmdouk.com/"
	url    = "http://www.mmdouk.com/forumdisplay.php?fid=59" // 网友自拍板块
	mkdir  = "D:/学习资料/"
	enc    = mahonia.NewEncoder("gbk") // 解码器
	client = http.Client{
		Timeout: 30 * time.Second,
	}
	syncMap map[string]string
	mysqlDB = stock.ActionMysql.Db
)

/**
爬取 Mimi!Board 论坛 网友自拍 图片
*/
func init() {
	RunWork()
}

func RunWork() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"),
	)
	c.SetRequestTimeout(20 * time.Second)

	ffHeader := func(request *colly.Request) {
		request.Headers.Add("Accept-Encoding", "gzip")
		request.Headers.Add("Cache-Control", "max-age=0")
		request.Headers.Add("Connection", "keep-alive")
		request.Headers.Add("Content-Type", "text/xml;charset=UTF-8")
		request.Headers.Add("Content-Encoding", "gzip")

		request.Headers.Add("Upgrade-Insecure-Requests", "1")
		request.Headers.Add("Host", "www.mmdouk.com")
		request.Headers.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		request.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	}

	imageColly, detailColly := c.Clone(), c.Clone()

	c.OnRequest(ffHeader)
	imageColly.OnRequest(ffHeader)
	c.OnResponse(func(response *colly.Response) {
		for i := 1; i <= 225; i++ {
			fmt.Printf("已成功连接到站点，正在获取数据，当前页：%d \n", i)
			imageColly.Visit(fmt.Sprintf("%s&page=%d", url, i))
		}
	})

	imageColly.OnResponse(func(response *colly.Response) {

	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err.Error())
	})

	// 获取分页 注：目前有效图片资源仅到 225 页！
	imageColly.OnHTML(".f_title > a", func(element *colly.HTMLElement) {
		var (
			title = until.ConvertToString(element.Text, "gbk", "utf-8")
			href  = base + element.Attr("href")
		)
		fmt.Println(title + "\n")

		cf := func() {
			if exist, _ := until.PathExists(mkdir + title); !exist {
				err := os.Mkdir(mkdir+title, os.ModePerm)
				if err != nil {
					print(err.Error())
				}
				ret, err := mysqlDB.Exec("INSERT image_list SET href = ?,path = ?", href, mkdir+title)
				if err != nil {
					panic(err)
				}
				if v, ok := ret.RowsAffected(); ok != nil || v <= 0 {
					log.Println("[写入数据库失败]:", ok.Error())
				}
			}
		}

		cf()

		detailColly.Visit(href)
	})

	detailColly.OnResponse(func(response *colly.Response) {
		var (
			l    = response.Request.URL.String()
			path = ""
		)
		_ = mysqlDB.QueryRow("select `path` from  image_list where href = ?", l).Scan(&path)
		response.Ctx.Put("mkdir", path)
	})

	detailColly.OnHTML(".t_msgfont > img", func(element *colly.HTMLElement) {
		var src = element.Attr("src")
		var path = element.Response.Ctx.Get("mkdir")
		request, mr := http.NewRequest("GET", src, nil)
		if mr != nil {
			fmt.Println("http error：" + mr.Error())
		}
		response, err := client.Do(request)
		if err != nil {
			log.Println(err.Error())
		} else {
			time := time.Now().Unix()
			var name = strconv.Itoa(int(time)) + until.RandSeq(4) + ".jpg"
			f, err := os.Create(path + "/" + name)
			if err == nil {
				content, _ := ioutil.ReadAll(response.Body)
				f.Write(content)
				f.Close()
			} else {
				fmt.Println("ERROR:" + err.Error())
			}
			fmt.Println(fmt.Sprintf("Write %s Succeed!", name))
		}
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Print(err)
	}
}
