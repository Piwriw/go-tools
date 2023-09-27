package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

/*
go get github.com/jordan-wright/email

邮件服务
*/

func EmailPool() {
	ch := make(chan *email.Email, 10)
	pool, err := email.NewPool(
		"smtp.126.com:25",
		4,
		smtp.PlainAuth("", "user@136.com", "password", "smtp.126.com"))
	if err != nil {
		logrus.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := pool.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}
	for i := 0; i < 10; i++ {
		e := email.NewEmail()
		e.From = "xx <yourname@136.com>"
		e.To = []string{"to@qq.com"}
		e.Subject = "Awesome web"
		e.Text = []byte(fmt.Sprintf("Awesome Web %d", i+1))
		ch <- e
	}

	close(ch)
	wg.Wait()
}

func main() {
	e := email.NewEmail()
	// 发送方
	e.From = "xxx"
	// 接收方
	e.To = []string{"toxxx"}
	// 邮件主题 标题
	e.Subject = "As"
	// 邮件正文
	e.Text = []byte("This is email text")
	// 抄送和 秘密抄送
	//e.Cc
	//e.Bcc
	// email 发送html美化
	//e.HTML
	// 添加附件
	_, err := e.AttachFile("go.mod")
	if err != nil {
		log.Fatalln(err)
	}
	err = e.Send("smtp.126.com:25", smtp.PlainAuth("", "xxx@126.com", "yyy", "smtp.126.com"))
	if err != nil {
		log.Fatal(err)
	}
}
