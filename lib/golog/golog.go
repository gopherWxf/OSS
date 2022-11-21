package golog

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

// 初始化方法自定义log
func init() {
	hostname, _ := os.Hostname()
	// 日志目录+日期.log
	file, err := os.OpenFile(fmt.Sprintf("%s%s.log", os.Getenv("LOG_DIRECTORY"), time.Now().Format("2006-01-02")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error myLog file: ", err)
	}
	//0.0.0.0-wxf [WARN] 2022/11/21 00:16:58 golog.go:37: I want a offer.
	//ip-hostname Level 2009/01/23 01:23:23 d.go:23 content

	Trace = log.New(file, os.Getenv("LISTEN_ADDRESS")+"-"+hostname+" [TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(file, os.Getenv("LISTEN_ADDRESS")+"-"+hostname+" [INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(file, os.Getenv("LISTEN_ADDRESS")+"-"+hostname+" [WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, os.Getenv("LISTEN_ADDRESS")+"-"+hostname+" [ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

//测试
func WriteLog() {
	for {
		Trace.Println("I want a offer.")
		Info.Println("I want a offer.")
		Warn.Println("I want a offer.")
		Error.Println("I want a offer.")
		time.Sleep(2 * time.Second)
	}
}
