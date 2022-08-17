package main

import (
	"ISP_Tool/config"
	"ISP_Tool/controller"
	"ISP_Tool/engine"
	"ISP_Tool/runningLog"
	"ISP_Tool/util"
	"ISP_Tool/view"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("软件开始运行。")
	err := controller.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败！")
	}
	ok := true
	for ok {
		newClient, err := util.Get_client()
		if err != nil {
			panic(err)
		}
		users, err := config.GetUserInfos()
		if err != nil {
			log.Println("GetUserInfos Error", err)
		}
		errConut := 0 //只通知一次
		for _, user := range users {
			err := engine.Run(newClient, user)
			if err != nil {
				fmt.Println()
				fmt.Println()
				log.Println("健康登记打卡失败,Error:", err)
				view.Clock_IN_Failed(user)
				fmt.Println()
				fmt.Println()
				errConut++
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
			fmt.Println("失败数量：", errConut)
			runningLog.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
		} else {
			runningLog.Inform("健康登记打卡 成功！")
		}
		fmt.Println()
		view.EndSlect()
		ok = controller.ProcessEndInput()
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
