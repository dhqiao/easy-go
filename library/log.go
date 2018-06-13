package library

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

// 日志处理，没有重写分级
func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	Logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

// 调试打印
func FmtPrint(split string, data ...interface{}) {
	fmt.Print(split)
	fmt.Print(data)
	fmt.Println(split)
}
