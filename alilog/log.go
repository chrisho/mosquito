package alilog

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	log "github.com/sirupsen/logrus"
)

// init log params
var (
	// 阿里云日志
	LogOff                 = strings.ToLower(os.Getenv("AliLogNoBoot")) == "true"
	LogStore               *sls.LogStore
	IpSource               string
	projectEndpoint        = os.Getenv("AliLogEndpoint")
	projectAccessKeyID     = os.Getenv("AliLogAccessKeyID")
	projectAccessKeySecret = os.Getenv("AliLogAccessKeySecret")
	projectName            = os.Getenv("AliLogName")
	logFile                = os.Getenv("AliLogFile")
	logStoreName           = os.Getenv("AliLogStoreName")
	LogTopic               = os.Getenv("AliLogTopic")
)

type Log struct {
	lineChan chan []string
}

func init() {
	// 不启动阿里云
	if LogOff || len(logStoreName) == 0 {
		return
	}
	newIpSource()
	newAliLog()

	// init log output
	fileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
	}
	mw := io.MultiWriter(os.Stdout, fileHandle)
	log.SetOutput(mw)

	NewLog()
}

// 本地ip
func newIpSource() {
	conn, err := net.Dial("tcp", "163.com:80")
	if err != nil {
		IpSource = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
}

// 阿里云客户端
func newAliLog() *sls.LogStore {
	// 配置
	logProject := &sls.LogProject{
		Name:            projectName,
		Endpoint:        projectEndpoint,
		AccessKeyID:     projectAccessKeyID,
		AccessKeySecret: projectAccessKeySecret,
	}
	// 实例化客户端
	var err error
	LogStore, err = logProject.GetLogStore(logStoreName)
	if err != nil {
		log.Error("logProject.GetLogStore error : " + err.Error())
	}
	return LogStore
}

func NewLog() {
	Log := &Log{
		lineChan: make(chan []string, 100),
	}
	Log.start()
}

func (Log *Log) start() {
	os.Remove(logFile)
	logf, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	mw := io.MultiWriter(logf, os.Stdout)
	log.SetOutput(mw)
	log.SetFormatter(&log.JSONFormatter{})

	go Log.getLogFile(Log.lineChan)
	go Log.pushToAliyun(Log.lineChan)
}

func (Log *Log) getLogFile(ch chan<- []string) {
	for {
		time.Sleep(10 * time.Second)

		logf, err := os.OpenFile(logFile, os.O_EXCL|os.O_RDONLY, 0666)
		if err != nil {
			log.Println(err)
		}
		reader := bufio.NewReader(logf)
		var readLines []string
		for {
			part, prefix, err := reader.ReadLine()
			if err != nil {
				break
			}
			if !prefix {
				readLines = append(readLines, string(part))
			}
		}
		if len(readLines) > 0 {
			ch <- readLines
		}
		os.Truncate(logFile, 0)
	}
}

func (Log *Log) pushToAliyun(ch <-chan []string) {
	// json struct
	type MessageLog struct {
		Level, Msg, Time string
	}
	var msgKey = "msg"
	var timeKey = "time"
	var levelKey = "level"
	// 信道
	for l := range ch {
		// 日志
		var slsLogs []*sls.Log
		// 监听
		for _, v := range l {
			var content []*sls.LogContent
			// json
			var message MessageLog
			json.Unmarshal([]byte(v), &message)
			// 日志内容
			content = append(content,
				&sls.LogContent{
					Key:   &msgKey,
					Value: &message.Msg,
				},
				&sls.LogContent{
					Key:   &timeKey,
					Value: &message.Time,
				},
				&sls.LogContent{
					Key:   &levelKey,
					Value: &message.Level,
				})
			// 日志内容
			if len(content) > 0 {
				timeNowUnix := uint32(time.Now().Unix())
				slsLogs = append(slsLogs, &sls.Log{
					Time:     &timeNowUnix,
					Contents: content,
				})
			}
		}
		// 发送日志
		if len(slsLogs) > 0 {
			logGroup := &sls.LogGroup{
				Topic:  &LogTopic,
				Source: &IpSource,
				Logs:   slsLogs,
			}
			if err := LogStore.PutLogs(logGroup); err != nil {
				log.Printf("PutLogs fail, err: %s\n", err)
			}
		}
	}
}
