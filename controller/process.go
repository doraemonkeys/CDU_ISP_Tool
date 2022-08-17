package controller

import (
	"ISP_Tool/config"
	"ISP_Tool/model"
	"ISP_Tool/view"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

//打卡结束处理用户的输入
func ProcessEndInput() bool {
	start := time.Now()
	ch := make(chan string)
	go endProcess(ch)
	var input string
	for time.Since(start) < time.Minute*2 {
		select {
		case input = <-ch:
			if input == "3" {
				return true
			}
			start = time.Now()
		default:
			time.Sleep(time.Second)
		}
	}
	return false
}

func endProcess(ch chan string) {
	var input string
	ok := true
	for ok {
		fmt.Scan(&input)
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)
		ch <- input
		switch input {
		case "0":
			err := config.DeleteUser()
			if err != nil {
				log.Println("删除失败！", err)
				fmt.Println("删除失败！", err)
			} else {
				fmt.Println("删除成功！")
			}
		case "1":
			err := config.ModifyUserInfos()
			if err != nil {
				if err.Error() == "取消修改" {
					break
				}
				log.Println("修改密码失败！", err)
				fmt.Println("修改密码失败！", err)
			} else {
				fmt.Println("修改密码成功！")
			}
		case "2":
			err := AddUser()
			if err != nil {
				fmt.Println("添加用户失败!")
			}
		case "3":
			ok = false
		case "4":
			err := config.RebuitConfig([]model.UserInfo{})
			if err != nil {
				fmt.Println("清空账号失败!")
				fmt.Println("请自己删除配置文件。")
			} else {
				fmt.Println("清空成功！")
			}
		}
		ch <- "keep_alive"
		fmt.Println(">>>>>>>>>>>无操作120秒后自动退出<<<<<<<<<<<")
		fmt.Println("请选择 【0 - 4】:")
	}
}

//添加用户信息
func AddUser() error {
	config, err := os.OpenFile("./配置文件.config", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("打开配置文件失败，Error:", err)
		fmt.Println("打开配置文件失败，Error:", err)
		return err
	}
	defer config.Close()
	var users []model.UserInfo
	var NewUser = model.UserInfo{}
	for {
		fmt.Println("输入 Q 退出添加账号")
		fmt.Println("请输入学号：")
		var id string
		fmt.Scan(&id)
		NewUser.UserID = strings.TrimSpace(id)
		if NewUser.UserID == "Q" || NewUser.UserID == "q" {
			break
		}
		fmt.Println("请输入密码：")
		var pwd string
		fmt.Scan(&pwd)
		NewUser.UserPwd = strings.TrimSpace(pwd)
		if NewUser.UserPwd == "Q" || NewUser.UserPwd == "q" {
			break
		}
		users = append(users, NewUser)
	}
	fmt.Scanf("\n")
	for _, v := range users {
		data, err := json.Marshal(v)
		if err != nil {
			log.Println("个人信息序列化失败！", err)
			fmt.Println("个人信息序列化失败！", err)
			return err
		}
		data = append(data, '\n')
		config.Write(data)
		fmt.Printf("添加 %s 成功！\n", v.UserID)
	}
	return nil
}

func InitConfig() error {
	log.Println("正在初始化配置文件")
	fmt.Println("正在初始化配置文件")
	fmt.Println()
	fmt.Println()
	view.Menu()
	fmt.Println()
	config, err := os.Open("./配置文件.config")
	if err == nil {
		defer config.Close()
		temp := make([]byte, 20)
		n, err := config.Read(temp)
		if err != nil && err != io.EOF {
			log.Println("预读取配置文件失败，Error:", err)
			fmt.Println("预读取配置文件失败，Error:", err)
			config.Close()
			return err
		}
		if n < 10 {
			log.Println("配置文件为空!")
			fmt.Println("配置文件为空!")
		} else {
			return nil
		}
	} else {
		log.Println("配置文件不存在或打开失败,Error:", err)
		fmt.Println("配置文件不存在或打开失败,Error:", err)
	}
	config, err = os.OpenFile("./配置文件.config", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("创建配置文件失败，Error:", err)
		fmt.Println("创建配置文件失败，Error:", err)
		return err
	}
	defer config.Close()
	fmt.Println()
	//用户协议
	ok := ReachConsensus()
	if !ok {
		panic("  用户不同意协议！")
	}
	fmt.Println("请输入你在成都大学的学号：")
	var id string
	fmt.Scan(&id)
	fmt.Println("请输入你的ISP密码：")
	var pwd string
	fmt.Scan(&pwd)
	fmt.Scanf("\n")
	var users []model.UserInfo
	var NewUser = model.UserInfo{}
	NewUser.UserID = strings.TrimSpace(id)
	NewUser.UserPwd = strings.TrimSpace(pwd)
	users = append(users, NewUser)
	for _, v := range users {
		data, err := json.Marshal(v)
		if err != nil {
			log.Println("个人信息序列化失败！", err)
			fmt.Println("个人信息序列化失败！", err)
			return err
		}
		data = append(data, '\n')
		config.Write(data)
	}
	fmt.Println()
	fmt.Println()
	return nil
}

func ReachConsensus() bool {
	fmt.Println(">>>>>>>>>>>>没有检测到账号，你可能是第一次使用本软件！<<<<<<<<<<<<<<<<")
	fmt.Println(">>>>>>>>>>>>          请阅读注意事项                  <<<<<<<<<<<<<<<<")
	fmt.Println()
	view.Warn()
	fmt.Println()
	var input string
	fmt.Println("是否同意？输入 [Yes] or [No]")
	fmt.Scan(&input)
	input = strings.TrimSpace(input)
	input = strings.ToUpper(input)
	for input != "YES" {
		if input == "NO" {
			return false
		}
		fmt.Scan(&input)
		fmt.Println("是否同意？输入 [Yes] or [No]")
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)
	}
	return true
}
