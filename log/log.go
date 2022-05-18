package log

import (
	"fmt"
	"os"
)

func Println(v ...interface{}) {
	println(fmt.Sprintln(v...))
}

func Printf(format string, v ...interface{}) {
	println(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	Println(v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	Printf(format, v...)
	os.Exit(1)
}
