package logx

import "fmt"

func Trace(i ...interface{}) {
	fmt.Println("TRACE ", i[0:])
}


func Tracef(format string, i ...interface{}) {
	fmt.Println("TRACE ", fmt.Sprintf(format, i...))
}


func DevDebugf(format string, i ...interface{}) {
	return
	fmt.Println("DevDEBUG ", fmt.Sprintf(format, i...))
}

func DevPrintf(format string, i ...interface{}) {
	fmt.Println("DevPrint ", fmt.Sprintf(format, i...))
}

func Debugf(format string, i ...interface{}) {
	fmt.Println("DEBUG ", fmt.Sprintf(format, i...))
}


func Debug(i ...interface{}) {
	fmt.Println("DEBUG ", i[0:])
}

func Warn(i ...interface{}){
	fmt.Println("WARN ", i[0:])
}

func Warnf(format string, i ...interface{}){
	fmt.Println("WARN ", fmt.Sprintf(format, i...))
}

func Error(i ...interface{}){
	fmt.Println("ERROR ", i[0:])
}

func Errorf(format string, i ...interface{}){
	fmt.Println("ERROR ", fmt.Sprintf(format, i...))
}