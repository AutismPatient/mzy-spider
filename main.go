package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	_ "mzy-spider/controller"
	"mzy-spider/httpreq"
	"mzy-spider/reg"
	"mzy-spider/stock"
	"mzy-spider/until"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const IsUPDATEREADY = true //是否更新密钥,仅用于调试模式

var Mode = true

func main() {

	if IsUPDATEREADY {
		token := until.GenerateToken(0)
		_, err := stock.Redis.Do("set", "run_key", token)
		if err != nil {
			panic(err)
		}
		// send email
		emailArr := strings.Split(httpreq.ConfigList["email"], ",")
		if len(emailArr) > 0 {
			err = until.SendEmail(emailArr, "用于启动的秘钥--"+until.TimeFormat(time.Now()), token, httpreq.ConfigList["授权码"])
			if err != nil {
				panic(err)
			}
		}
	}
	// 判断是否启用快捷模式
	actionMysql := stock.ActionMysql.Db
	dateLine := int64(0)

	err := actionMysql.QueryRow("SELECT dateline FROM movie_info ORDER BY dateline DESC LIMIT 1").Scan(&dateLine)
	if err != nil {
		panic(err)
	}

	day, _ := strconv.ParseInt(httpreq.ConfigList["快捷模式"], 0, 64)
	if (time.Now().Unix() - dateLine) < (86400 * day) { // 5天内启用
		_, err := stock.Redis.Do("set", "next", time.Now().Unix())
		if err != nil {
			panic(err)
		}
		log.Println(until.TimeFormat(time.Now()) + "快捷模式")
	}

	until.PrintlnMsg(false, true, until.TimeFormat(time.Now())+" 站点初始化成功，秘钥已更新")

	server := gin.New() // New()
	//set Mode
	if Mode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// custom format
		var statusColor, methodColor, resetColor, byteStr string
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
		var offset = reg.Decimal(float64(param.BodySize) / 1024)
		if offset <= 0 {
			byteStr = "0"
		} else {
			byteStr = strconv.FormatFloat(offset, 'f', -1, 32)
		}
		return fmt.Sprintf("%s - [%s] - [%s %s %s %13v %s %3d %s %s %s] - [%s kb] - [%s]\n",
			param.ClientIP,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			methodColor, param.Method, resetColor,
			param.Latency,
			statusColor, param.StatusCode, resetColor,
			param.Request.Proto,
			param.Path,
			byteStr,
			param.Request.UserAgent(),
		)
	}))
	server.Use(gin.Recovery()) // 消除异常
	//注册路由到 gin handler
	for path, handle := range reg.Actions {
		if v, ok := reg.Ignore[path]; ok {
			server.GET(path, v)
			server.POST(path, v)
			continue
		}
		if v, ok := reg.PathActions[path]; ok {
			group := server.Group(v)
			if auth, ok := reg.AuthFunc[v]; ok {
				group.Use(reg.ServerMiddleware(auth))
			}
			if strings.Contains(path, v) {
				path = strings.Replace(path, v, "", -1)
			}
			//set get、post httpMethod
			group.GET(path, handle)
			group.POST(path, handle)
		}
	}
	server.Static("/static", "html")

	srv := &http.Server{
		Addr:           ":8888",
		Handler:        server,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    8 * time.Second,
	}
	srv.ListenAndServe()
}
