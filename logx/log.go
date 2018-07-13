package logx

import (
	"fmt"
	"time"
)

func Trace(i ...interface{}) {
	fmt.Println(getLogTime() + " TRACE ", i[0:])
}


func Tracef(format string, i ...interface{}) {
	fmt.Println(getLogTime() + " TRACE ", fmt.Sprintf(format, i...))
}


func DevDebugf(format string, i ...interface{}) {
	fmt.Println(getLogTime() + " DevDEBUG ", fmt.Sprintf(format, i...))
}

func DevPrintf(format string, i ...interface{}) {
	fmt.Println(getLogTime() + " DevPrint ", fmt.Sprintf(format, i...))
}

func Debugf(format string, i ...interface{}) {
	fmt.Println(getLogTime() + " DEBUG ", fmt.Sprintf(format, i...))
}


func Debug(i ...interface{}) {
	fmt.Println(getLogTime() + " DEBUG ", fmt.Sprint(i...))
}

func Info(i ...interface{}) {
	fmt.Println(getLogTime() + " INFO ", fmt.Sprint(i...))
}

func Infof(format string, i ...interface{}){
	fmt.Println(getLogTime() + " INFO ", fmt.Sprintf(format, i...))
}

func Warn(i ...interface{}){
	fmt.Println(getLogTime() + " WARN ", fmt.Sprint(i...))
}

func Warnf(format string, i ...interface{}){
	fmt.Println(getLogTime() + " WARN ", fmt.Sprintf(format, i...))
}

func Error(i ...interface{}){
	fmt.Println(getLogTime() + " ERROR ", fmt.Sprint(i...))
}

func Errorf(format string, i ...interface{}){
	fmt.Println(getLogTime() + " ERROR ", fmt.Sprintf(format, i...))
}

func getLogTime() string{
	return time.Now().Format("2006-01-02 15:04:05")
}