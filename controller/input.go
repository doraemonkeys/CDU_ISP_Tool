package controller

import (
	"ISP_Tool/model"
	"ISP_Tool/server"
	"ISP_Tool/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

// 打卡结束处理用户的输入,返回是否继续循环
func ProcessEndInput() bool {
	var input string
	for {
		fmt.Scan(&input)
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)

		switch input {
		case "0":
			err := server.DeleteUser()
			if err != nil {
				log.Println("删除失败！", err)
				fmt.Println("删除失败！", err)
			} else {
				fmt.Println("删除成功！")
				model.UserConfigChanged = true
			}
		case "1":
			err := server.ModifyUserInfos()
			if err != nil {
				if err.Error() == "取消修改" {
					break
				}
				log.Println("修改密码失败！", err)
				color.Red("修改密码失败！%s", err.Error())
				//fmt.Println("修改密码失败！", err)
			} else {
				color.HiGreen("修改密码成功！")
				//fmt.Println("修改密码成功！")
				model.UserConfigChanged = true
			}
		case "2":
			model.UserConfigChanged = true
			err := server.AddUser()
			if err != nil {
				fmt.Println("添加用户失败!")
			}
		case "3":
			return true
		case "4":
			err := server.RebuildConfig([]model.UserInfo{})
			if err != nil {
				fmt.Println("清空账号失败!")
				fmt.Println("请自己删除配置文件。")
			} else {
				fmt.Println("清空成功！")
				model.UserConfigChanged = true
			}
		case "5":
			fmt.Println()
			if model.Auto_Start {
				err := server.CancelAutoStart()
				if err != nil {
					log.Println("关闭自启动失败！", err)
					fmt.Println("关闭自启动失败！", err)
				} else {
					color.HiGreen("关闭自启动成功！")
					log.Println("关闭自启动成功！")
					model.Auto_Start = false
				}
			} else {
				err := server.SetAutoStart()
				if err != nil {
					log.Println("开启自启动失败！")
					fmt.Println("开启自启动失败！")
				} else {
					fmt.Println()
					color.HiGreen("开启自启动成功！ 程序正在重新启动.....")
					log.Println("开启自启动成功！")
					model.Auto_Start = true
					//启动一个自动打卡守护进程，当前进程默认退出
					os.Exit(0)
				}
			}
			fmt.Println()
		case "6":
			err := server.SwitchChooseLocation()
			if err != nil {
				log.Println("切换地址选择方式失败！", err)
				color.Red("切换地址选择方式失败！%s", err.Error())
			} else {
				color.HiGreen("切换地址选择方式成功！\n")
				log.Println("切换地址选择方式成功！")
			}
		}
		fmt.Println()
		attributes := [4]color.Attribute{}
		attributes[1] = color.FgRed
		utils.ColorPrint(attributes[:], "请选择 【", "0 - 6", "】:\n")
	}
}
