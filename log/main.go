package logs

import (
	"log"
	"os"
	"time"
)

var loger *log.Logger

func InitLog() {

	file := "./log/" + time.Now().Format("2006-01-02-15") + "_log" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为loger作为输出

}
