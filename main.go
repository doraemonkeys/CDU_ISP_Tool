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
	//检查更新
	server.CheckUpdate()
	shouldReturn := checkStatus()
	if shouldReturn {
		return
	}
	//for用于更新配置文件后可能的重新打卡
	for {
		users, err := server.GetUserInfos()
		if err != nil {
			log.Println("GetUserInfos Error", err)
			fmt.Println("GetUserInfos Error", err)
		}
		//打卡,返回打卡失败的数量
		errCount := checkIn(users)
		//不管打卡成功或者失败，循环一次后都认为配置文件已经执行过一次,用户没有再次修改
		model.UserConfigChanged = false

		//处理打卡结果
		workWithCheckInResults(errCount, users)
		fmt.Println()

		showAutoStartStatus()
		fmt.Println()

		view.EndSelect()
		ok := controller.ProcessEndInput()
		if !ok {
			return
		}
	}
}

func showAutoStartStatus() {
	fmt.Printf("自动打卡状态：")
	if !model.Auto_Start {
		color.Red("已关闭")
	} else {
		color.HiGreen("已开启")
	}
}

func workWithCheckInResults(errCount int, users []model.UserInfo) {
	if errCount != 0 {
		color.Red("打卡失败数量：%d", errCount)
		log.Println("打卡失败数量：", errCount)
		if model.Auto_Start && !model.Auto_Clock_IN_Success {
			server.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")

			server.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡失败")
			model.Auto_Clock_IN_Success = false
		}
		if !model.Auto_Start {
			server.Inform("健康登记打卡 失败！！！" + " 请手动打卡。")
		}
	}
	if errCount == 0 && len(users) != 0 {
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
}

func checkIn(users []model.UserInfo) int {
	errCount := 0
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
			errCount++
		} else {
			fmt.Println()
			fmt.Println()
			log.Println(user.UserID, "健康登记打卡成功")
			view.Clock_IN_Success(user)
			fmt.Println()
			fmt.Println()
		}
	}
	return errCount
}

func checkStatus() bool {
	err := server.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败！")
		if server.CheckAutoStart() && !server.TodayCheckInSuccess() {
			server.WrittenToTheLog(time.Now().Format("2006/01/02") + " 自动打卡失败")
		}
		time.Sleep(time.Minute)
		return true
	}
	return false
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
