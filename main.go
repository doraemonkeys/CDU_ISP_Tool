package main

import (
	"ISP_Tool/controller"
	"ISP_Tool/engine"
	"ISP_Tool/model"
	"ISP_Tool/server"
	"ISP_Tool/utils"
	"ISP_Tool/view"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

func main() {
	log.Println()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>软件开始运行。")
	err := server.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败！")
		time.Sleep(time.Minute)
		return
	}
	//for用于更新配置文件后可能的重新打卡
	for {
		users, err := server.GetUserInfos()
		if err != nil {
			log.Println("GetUserInfos Error", err)
			fmt.Println("GetUserInfos Error", err)
		}
		errConut := 0 //打卡失败数量
		for _, user := range users {
			newClient, _ := utils.Get_client()
			err := engine.Run(newClient, user)
			if err != nil && err.Error() != "健康登记打卡已存在" {
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
		//不管打卡成功或者失败，循环一次后都认为配置文件已经执行过一次,用户没有再次修改
		model.UserConfigChanged = false
		if errConut != 0 {
			color.Red("打卡失败数量：%d", errConut)
			log.Println("打卡失败数量：", errConut)
			if model.Auto_Start && !model.Auto_Clock_IN_Success {
				server.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
				//写入自动打卡日志
				server.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡失败")
				model.Auto_Clock_IN_Success = false
			}
			if !model.Auto_Start {
				server.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
			}
		}
		if errConut == 0 && len(users) != 0 {
			if model.Auto_Start && !model.Auto_Clock_IN_Success {
				server.Inform("健康登记打卡 成功！")
				server.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡成功")
				model.Auto_Clock_IN_Success = true
			}
			if !model.Auto_Start {
				server.Inform("健康登记打卡 成功！")
			}
		}
		if len(users) == 0 {
			fmt.Println()
			fmt.Println("在配置文件中没有找到用户！")
			log.Println("在配置文件中没有找到用户！")
		}
		fmt.Println()
		fmt.Printf("自动打卡状态：")
		if !model.Auto_Start {
			color.Red("已关闭")
		} else {
			color.HiGreen("已开启")
		}
		fmt.Println()
		view.EndSlect()
		ok := controller.ProcessEndInput()
		if !ok {
			return
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
