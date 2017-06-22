package alilog

import (
	"os"
	"io"
	"net"
	"bufio"
	"time"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aliyun/aliyun-log-go-sdk/example/util"
	"github.com/gogo/protobuf/proto"
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

// init log params
var (
	projectEndpoint = os.Getenv("AliLogEndpoint")
	projectAccessKeyID = os.Getenv("AliLogAccessKeyID")
	projectAccessKeySecret = os.Getenv("AliLogAccessKeySecret")
	projectName = os.Getenv("AliLogName")
	logFile = os.Getenv("AliLogFile")
	logStoreName = os.Getenv("AliLogStoreName")
	logTopic = os.Getenv("AliLogTopic")
)

type Log struct {
	LOG_FILE string

	LOGSTORE_NAME string

	LOG_TOPIC  string
	LOG_SOURCE string

	lineChan chan []string
}

func CatchPanic() {
	if err := recover(); err != nil {
		log.Error(err)
	}
}

func NewLog() {

	conn, err := net.Dial("tcp", "163.com:80")
	ip := ""
	if err == nil {
		defer conn.Close()
		ip = strings.Split(conn.LocalAddr().String(), ":")[0]
	}

	Log := &Log{
		logFile,
		logStoreName,
		logTopic,
		ip,
		make(chan []string, 100),
	}

	util.Project.Endpoint = projectEndpoint
	util.Project.AccessKeyID = projectAccessKeyID
	util.Project.AccessKeySecret = projectAccessKeySecret
	util.Project.Name = projectName

	Log.Start()
}

func (Log *Log) Start() {
	os.Remove(Log.LOG_FILE)
	logf, err := os.OpenFile(Log.LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	mw := io.MultiWriter(logf, os.Stdout)
	log.SetOutput(mw)
	log.SetFormatter(&log.JSONFormatter{})

	go Log.GetLogFile(Log.lineChan)
	go Log.PushToAliyun(Log.lineChan)
}

func (Log *Log) GetLogFile(ch chan<- []string) {
	for {
		time.Sleep(10 * time.Second)

		logf, err := os.OpenFile(Log.LOG_FILE, os.O_EXCL|os.O_RDONLY, 0666)
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
		os.Truncate(Log.LOG_FILE, 0)
	}
}

func (Log *Log) PushToAliyun(ch <-chan []string) {
	logstore_name := Log.LOGSTORE_NAME

	logstore, err := util.Project.GetLogStore(logstore_name)
	if err != nil {
		log.Printf("GetLogStore fail, err: ", err)
		return
	}

	type MessageLog struct {
		Level, Msg, Time string
	}

	for l := range ch {

		slslogs := []*sls.Log{}
		for _, v := range l {
			content := []*sls.LogContent{}

			var message MessageLog
			json.Unmarshal([]byte(v), &message)

			content = append(content,
				&sls.LogContent{
					Key:   proto.String("msg"),
					Value: proto.String(message.Msg),
				},
				&sls.LogContent{
					Key:   proto.String("time"),
					Value: proto.String(message.Time),
				},
				&sls.LogContent{
					Key:   proto.String("level"),
					Value: proto.String(message.Level),
				})

			if content != nil {
				slslogs = append(slslogs, &sls.Log{
					Time:     proto.Uint32(uint32(time.Now().Unix())),
					Contents: content,
				})
			}
		}

		if slslogs != nil {
			loggroup := &sls.LogGroup{
				Topic:  proto.String(Log.LOG_TOPIC),
				Source: proto.String(Log.LOG_SOURCE),
				Logs:   slslogs,
			}
			err = logstore.PutLogs(loggroup)
			if err != nil {
				log.Printf("PutLogs fail, err: %s\n", err)
			}
		}
	}
}
