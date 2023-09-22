package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var Log = log.New()

func init() {
	Log = NewLogger()
}

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func NewLogger() *log.Logger {
	filepaths := "./log/sys.log"
	writer, _ := rotatelogs.New(
		filepaths+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(filepaths),
		rotatelogs.WithMaxAge(time.Duration(604800)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  writer,
		log.FatalLevel: writer,
		log.DebugLevel: writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.PanicLevel: writer,
	}
	Log.SetReportCaller(true)
	lfHook := lfshook.NewHook(writeMap, &MyFormatter{})
	Log.AddHook(lfHook)
	return Log
}
