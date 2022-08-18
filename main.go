package main

import (
	"ISP_Tool/config"
	"ISP_Tool/controller"
	"ISP_Tool/engine"
	"ISP_Tool/model"
	"ISP_Tool/runningLog"
	"ISP_Tool/util"
	"ISP_Tool/view"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.Println()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>软件开始运行。")
	err := controller.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败！")
	}
	ok := true
	for ok {
		newClient, err := util.Get_client()
		if err != nil {
			log.Println("创建虚拟客户端失败！", err)
			time.Sleep(time.Hour / 2)
			continue
		}
		users, err := config.GetUserInfos()
		if err != nil {
			log.Println("GetUserInfos Error", err)
		}
		errConut := 0 //只通知一次
		for _, user := range users {
			err := engine.Run(newClient, user)
			if err != nil {
				if err.Error() != "健康登记打卡已存在" {
					fmt.Println()
					fmt.Println()
					log.Println("健康登记打卡失败,Error:", err)
					view.Clock_IN_Failed(user)
					fmt.Println()
					fmt.Println()
					errConut++
				}
			} else {
				fmt.Println()
				fmt.Println()
				log.Println(user.UserID, "健康登记打卡成功")
				view.Clock_IN_Success(user)
				fmt.Println()
				fmt.Println()
			}
		}
		if errConut != 0 {
			fmt.Println("打卡失败数量：", errConut)
			log.Println("打卡失败数量：", errConut)
			if model.Auto_Start {
				if !model.Auto_Clock_IN_Success {
					runningLog.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
					controller.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡失败")
				}
			} else {
				runningLog.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
			}
		}
		if errConut == 0 && len(users) != 0 {
			if model.Auto_Start {
				if !model.Auto_Clock_IN_Success {
					runningLog.Inform("健康登记打卡 成功！")
					controller.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡成功")
				}
			} else {
				runningLog.Inform("健康登记打卡 成功！")
			}
		}
		if len(users) == 0 {
			fmt.Println()
			fmt.Println("在配置文件中没有找到用户！")
			log.Println("在配置文件中没有找到用户！")
		}
		fmt.Println()
		view.EndSlect()
		ok = controller.ProcessEndInput()
		if !ok && model.Auto_Start {
			sleepTime := time.Now().Format("2006/01/02")
			for sleepTime == time.Now().Format("2006/01/02") {
				log.Println("程序休眠一小时")
				time.Sleep(time.Hour)
			}
			ok = true //第二天到了
		}
	}
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
