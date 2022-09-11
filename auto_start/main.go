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
		//获取今日打卡失败次数
		failCount, err := server.FailedTimes()
		if err != nil {
			log.Println("获取今日打卡失败次数失败，Error:", err)
			time.Sleep(time.Minute * 2)
			continue
		}
		//第二天到了后,12:10分之后才开始打卡
		if time.Now().Hour() == 0 && time.Now().Minute() < 10 {
			time.Sleep(time.Minute)
			continue
		}
		//检查今日自动打卡是否成功,今日自动打卡是否存在失败记录(失败2次则不启动打卡程序),
		if !server.TodayClockInSuccess() {
			//如果第一次打卡失败，第二次打卡时间为早上7:00之后
			if (failCount == 0) || (failCount <= 2 && time.Now().Hour() >= 7) {
				//检查更新
				server.CheckUpdate()
				log.Println("检测到自动打卡未执行,正在启动打卡程序。")
				err := server.StartNewProgram()
				if err != nil {
					log.Println("启动打卡程序主体失败！", err)
				}
			}
		}
		//time.Sleep(time.Hour / 2)
		time.Sleep(time.Minute * 2)
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
