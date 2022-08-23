package server

import (
	"log"
	"os"
)

func WrittenToTheLog(content string) {
	CheckAutoStart()
	auto_start, err := os.OpenFile("./config/auto_start.config", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("自动打卡写入日志失败！", err)
		return
	}
	defer auto_start.Close()
	_, err = auto_start.WriteString(content + "\n")
	if err != nil {
		log.Println("自动打卡写入日志失败！", err)
		return
	}
}
