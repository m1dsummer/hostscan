package elog

import (
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"time"
)

var format logging.Formatter

func init() {
	format = logging.MustStringFormatter(`%{color}[%{level:.4s}]%{color:reset} [%{time:15:04:05}] %{message}`)
}

type LogFiller struct {
	//record *logging.Record
	msg string
}

func NewFiller(record *logging.Record) *LogFiller {
	return &LogFiller{msg: toMsg(record)}

}
func toMsg(record *logging.Record) string {
	buf := new(bytes.Buffer)
	format.Format(4, record, buf)
	return buf.String()
}

func LogWithLevel(msg string, level logging.Level) {
	record := &logging.Record{
		Time: time.Now(),
		//Module: "",
		Args:  []interface{}{msg},
		Level: level,
	}
	//println(msg)
	fmt.Println(toMsg(record))
}

func Info(msg string) {
	LogWithLevel(msg, logging.INFO)
}

func Debug(msg string) {
	LogWithLevel(msg, logging.DEBUG)
}

func Warn(msg string) {
	LogWithLevel(msg, logging.WARNING)
}

func Error(msg string) {
	LogWithLevel(msg, logging.ERROR)
}

func Notice(msg string) {
	LogWithLevel(msg, logging.NOTICE)
}

func Critical(msg string) {
	LogWithLevel(msg, logging.CRITICAL)
}
