package until

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
	"time"
)

const (
	defaultKey = "HSJnpYbnLamBhu"
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
func SendEmail(mailTo []string, subject string, body string) error {

	mailConn := map[string]string{
		"user": "1010014622@qq.com",
		"pass": "tvyzsstspscrbbjh", // 授权码
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()
	m.SetHeader("From", "程序启动秘钥"+"<"+mailConn["user"]+">")
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}
