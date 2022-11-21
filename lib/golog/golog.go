package golog

import (
	es "OSS/lib/ElasticSearch"
	"fmt"
	"github.com/nxadm/tail"
	"log"
	"os"
	"strconv"
	"strings"
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

// 字符串时间转时间戳工具
func dateFormat(dataTime string) int64 {
	loc, _ := time.LoadLocation("Local")
	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", dataTime, loc)

	if err != nil {
		Error.Printf("filed datetime format")
	}
	return formatTime.Unix()
}

// ReadLog 实时读取日志文件并推送到ES
func ReadLog(times string) {
	fileName := fmt.Sprintf("%s%s.log", os.Getenv("LOG_DIRECTORY"), times)
	config := tail.Config{
		//const (
		//	SeekStart   = 0 // seek relative to the origin of the file
		//	SeekCurrent = 1 // seek relative to the current offset
		//	SeekEnd     = 2 // seek relative to the end
		//)
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // Location读取文件的位置, Whence=2=SeekEnd
		ReOpen:    true,                                 // tail -F
		Follow:    true,                                 // tail -f
		MustExist: false,                                // 如果文件不存在，则提前失败--->可以不存在
		Poll:      true,                                 // 轮询
	}
	// 打开文件读取日志
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		Error.Println("tail file failed, err:", err)
		return
	}
	// 开始读取数据
	for {
		msg, ok := <-tails.Lines
		if !ok {
			Error.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second) // 读取出错停止一秒
			continue
		}
		//切割成json
		docMsg := analysisLog(msg.Text)
		// 异步向es推送log
		go es.PutLog(docMsg)
	}
}

// 解析某条log
func analysisLog(logContent string) string {
	// TODO host_level_date_time_content
	// {"osName":"%s","level":"%s","dateTime":%d,"content":"%s"}`

	//0.0.0.0-wxf [WARN] 2022/11/21 00:16:58 golog.go:37: I want a offer.
	//ip-hostname Level 2009/01/23 01:23:23 d.go:23 content

	//[0.0.0.0-wxf,[WARN],2022/11/21,00:16:58,golog.go:37:,=====>many--->"I want a offer."]
	splitstring := strings.Split(logContent, " ")

	osName := splitstring[0]
	level := strings.Trim(splitstring[1], "[]")
	date := strings.Replace(splitstring[2], "/", "-", -1)
	time0 := splitstring[3]
	//abc123aa   3+3
	sidx := strings.Index(logContent, time0) + len(time0)
	//返回一个双引号的Go字符串文本
	content := strconv.Quote(logContent[sidx:])
	doc := fmt.Sprintf(`{"osName":"%s","level":"%s","dateTime":%d,"content":%s}`,
		osName, level, dateFormat(fmt.Sprintf("%s %s", date, time0)), content)
	return doc
}
