package until

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultKey = "HSJnpYbnLamBhu"
	C          = "C:/视频资源"
	D          = "D:/视频资源"
	E          = "E:/视频资源"
	F          = "F:/视频资源"
)

func GenerateToken(offset int64) string {
	var time = time.Now().Unix()
	if offset != 0 {
		time += offset
	}
	str := []byte(fmt.Sprintf("%s%d%s", defaultKey, time, defaultKey))
	newToken := md5.Sum(str)
	return fmt.Sprintf("%X", newToken)
}
func PrintlnMsg(error, ln bool, msg string) {
	str := "INFO"
	if error {
		str = "ERROR"
	}
	str = fmt.Sprintf("[SYSTEM %s]:", str) + msg
	if ln {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

// 发送邮件
func SendEmail(mailTo []string, subject, body, pass string) error {

	mailConn := map[string]string{
		"user": "1010014622@qq.com",
		"pass": pass, // 授权码
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()
	m.SetHeader("From", "RUN KEY"+"<"+mailConn["user"]+">")
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}
func Json(resp http.ResponseWriter, data interface{}) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Content-Type", "application/json")

	resp.Write(buffer.Bytes())
}
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
