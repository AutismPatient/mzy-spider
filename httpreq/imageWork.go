package httpreq

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
	"github.com/gomodule/redigo/redis"
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
	base            = "http://www.mmdouk.com/"
	url             = "http://www.mmdouk.com/forumdisplay.php?fid=59" // 网友自拍板块
	page            = 230
	mkdir           = "D:/学习资料/"
	NilMkdir        = "D:/学习资料/散"
	PageRedisMkdir  = "page_href_list:"
	ImageRedisMkdir = "images_list:"
	enc             = mahonia.NewEncoder("gbk") // 解码器
	client          = http.Client{
		Timeout: 6 * time.Second,
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
		for i := 1; i <= page; i++ { // todo 多线程 2021年6月25日14:12:14
			f := func(i int) {
				fmt.Printf("已成功连接到站点，正在获取数据，当前页：%d \n", i)
				err := imageColly.Visit(fmt.Sprintf("%s&page=%d", url, i))
				if err != nil {
					return
				}
			}
			f(i)
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
			skip  = false
		)
		if title == "" || href == "" {
			return
		}

		fmt.Println(title + "\n")

		_, err := redis.String(stock.Redis.Do("get", PageRedisMkdir+href))

		if err != nil && err == redis.ErrNil {
			skip = true
		}

		cf := func() {
			if exist, _ := until.PathExists(mkdir + title); !exist {
				err := os.Mkdir(mkdir+title, os.ModePerm)
				_, err = stock.Redis.Do("set", PageRedisMkdir+href, mkdir+title)

				if err != nil {
					print(err.Error())
				}
			}
		}

		cf()

		if skip {
			err = detailColly.Visit(href)
			if err != nil {
				return
			}
		}
	})

	detailColly.OnResponse(func(response *colly.Response) {
		//var (
		//	l    = response.Request.URL.String()
		//	path = ""
		//)
		//_ = mysqlDB.QueryRow("select `path` from  image_list where href = ?", l).Scan(&path)
		//response.Ctx.Put("mkdir", path)
	})

	detailColly.OnHTML(".t_msgfont > img", func(element *colly.HTMLElement) {
		var (
			src  = element.Attr("src")
			path = element.Response.Request.URL.String() // todo 2021年6月25日14:10:38
			unix = time.Now().Unix()
			name = strconv.Itoa(int(unix)) + until.RandSeq(4) + ".jpg"

			do = func(name string) {
				_, err := stock.Redis.Do("set", ImageRedisMkdir+src, name)

				if err != nil {
					print(err.Error())
				}
			}
		)

		mkdir, err := redis.String(stock.Redis.Do("get", PageRedisMkdir+path))

		if err != nil && err == redis.ErrNil {
			mkdir = NilMkdir
		}

		_, err = redis.String(stock.Redis.Do("get", ImageRedisMkdir+src))

		if err == nil {
			return
		}

		request, mr := http.NewRequest("GET", src, nil)
		if mr != nil {
			fmt.Println("http error：" + mr.Error())
			return
		}
		response, err := client.Do(request)
		if err != nil {
			log.Println(err.Error())
		} else {
			f, err := os.Create(mkdir + "/" + name)
			if err == nil {
				content, _ := ioutil.ReadAll(response.Body)
				_, err := f.Write(content)
				if err != nil {
					return
				}
				err = f.Close()
				if err != nil {
					return
				}
			} else {
				fmt.Println("ERROR:" + err.Error())
			}

			err = response.Body.Close()
			if err != nil {
				return
			}

			fmt.Println(fmt.Sprintf("Write %s Succeed!", name))
		}

		do(name)

	})

	err := c.Visit(url)
	if err != nil {
		fmt.Print(err)
	}
}
