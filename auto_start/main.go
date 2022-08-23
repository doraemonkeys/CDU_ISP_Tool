package main

import (
	"fmt"
	"ispTool_auto_start/server"
	"log"
	"os"
	"time"
)

//开机自启动或者第一次开启自启动时，本程序将于后台运行，完成自动打开的功能。
//所有配置文件都放在config文件夹
func main() {
	log.Println(">>>>>>>>>>>>>>>>>>>>>自动打卡守护程序开始运行。")
	for {
		//自动打卡是否开启
		if !server.CheckAutoStart() {
			log.Println("检测到自动打卡已关闭,程序即将退出。")
			break
		}
		//检查用户信息配置文件是否存在
		server.ConfigFileExist()
		//检查今日自动打卡是否成功
		if !server.TodayClockInSuccess() {
			log.Println("检测到自动打卡未成功,正在启动打卡程序。")
			err := server.StartNewProgram()
			if err != nil {
				log.Println("启动打卡程序主体失败！", err)
			}
		}
		log.Println("自动打卡守护程序休眠半小时")
		time.Sleep(time.Hour / 2)
	}
	log.Println(">>>>>>>>>>>>>>>>>>>>>自动打卡守护程序退出。")
}

func init() {
	logFile, err := os.OpenFile("RunningLog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Ldate)
}
