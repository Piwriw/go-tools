package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

/*
https://github.com/sirupsen/logrus
go get github.com/sirupsen/logrus
logrus完全兼容标准的log库，还支持文本、JSON 两种日志输出格式
*/
func main() {

	logrus.SetLevel(logrus.TraceLevel)
	//
	//logrus.Trace("trace msg")
	//logrus.Debug("debug msg")
	//logrus.Info("info msg")
	//logrus.Warn("warn msg")
	//logrus.Error("error msg")
	//logrus.Fatal("fatal msg")
	//logrus.Panic("panic msg")
	// 添加输出字段
	logerLog := logrus.WithFields(logrus.Fields{
		"name": "piwriw",
		"age":  18,
	})
	logerLog.Info("info with fields")
	file, err := os.OpenFile("./logrus/logrus.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.SetOutput(file)
	logrus.Info("info msg")

}
