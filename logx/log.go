package logx

import "fmt"

func Tracef(format string, i ...interface{}) {
	fmt.Println("TRACE ", fmt.Sprintf(format, i...))
}


func Debugf(format string, i ...interface{}) {
	fmt.Println("DEBUG ", fmt.Sprintf(format, i...))
}
