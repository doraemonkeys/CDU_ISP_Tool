package controller

import (
	"ISP_Tool/config"
	"ISP_Tool/model"
	"ISP_Tool/util"
	"ISP_Tool/view"
	"bufio"
	"encoding/json"
	"errors"
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
	//2分钟无操作返回false,如果没有设置自启动,程序会退出。
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
		case5Count := 3 //case 5限制最多执行3次
		switch input {
		case "0":
			err := config.DeleteUser()
			if err != nil {
				log.Println("删除失败！", err)
				fmt.Println("删除失败！", err)
			} else {
				fmt.Println("删除成功！")
				model.UserConfigChanged = true
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
				model.UserConfigChanged = true
			}
		case "2":
			model.UserConfigChanged = true
			err := config.AddUser()
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
				model.UserConfigChanged = true
			}
		case "5":
			if case5Count == 0 {
				fmt.Println("设置更改次数过多,请重启电脑保证程序正确运行！")
				log.Println("设置更改次数过多,请重启电脑保证程序正确运行！")
				break
			}
			case5Count--
			fmt.Println()
			fmt.Println(">>>>>如果设置失败，请关闭杀毒软件并以管理员权限重新运行<<<<<")
			fmt.Println()
			if model.Auto_Start {
				err := config.CancelAutoStart()
				if err != nil {
					log.Println("关闭自启动失败！")
					fmt.Println("关闭自启动失败！")
				} else {
					fmt.Println("关闭自启动成功！")
					log.Println("关闭自启动成功！")
					model.Auto_Start = false
				}
			} else {
				err := config.SetAutoStart()
				if err != nil {
					log.Println("开启自启动失败！")
					fmt.Println("开启自启动失败！")
				} else {
					fmt.Println("开启自启动成功！")
					log.Println("开启自启动成功！")
					model.Auto_Start = true
				}
			}
			fmt.Println()
		}
		ch <- "keep_alive"
		fmt.Println(">>>>>>>>>>>无操作120秒后自动退出<<<<<<<<<<<")
		fmt.Println("请选择 【0 - 5】:")
	}
}

func InitConfig() error {
	log.Println("正在初始化配置文件")
	fmt.Println("正在初始化配置文件")
	model.Auto_Start = CheckAutoStart()
	if model.Auto_Start {
		if TodayClockInSuccess() {
			model.Auto_Clock_IN_Success = true
			view.Auto_Clock_IN_Success()
			fmt.Println()
			startTime := time.Now()
			fmt.Printf("按Enter键继续执行程序......")
			ch := make(chan bool, 1)
			go util.PressToContinue(ch)
			ok := false
			for time.Since(startTime) < time.Minute/2 {
				select {
				case ok = <-ch:
				default:
					time.Sleep(time.Second / 4)
				}
				if ok {
					break
				}
			}
		}
	}
	fmt.Println("正在检查网络环境...")
	for i := 0; !util.NetWorkStatus(); i++ {
		time.Sleep(time.Second)
		if i == 30 {
			fmt.Println("网络连接错误，请检查网络配置!")
			log.Println("网络连接错误，请检查网络配置!")
			return errors.New("网络连接错误")
		}
	}
	fmt.Println("Net Status , OK!")
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
	model.UserConfigChanged = true
	return nil
}

//是否已经设置为自启动
func CheckAutoStart() bool {
	autoStart, err := os.Open("./auto_start.config")
	if err != nil {
		log.Println("读取自启动信息失败！", err)
		autoStart, err := os.OpenFile("./auto_start.config", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			log.Println("初始化自启动信息失败！", err)
		}
		defer autoStart.Close()
		//默认设置为不自启动
		autoStart.WriteString("false\n")
		return false
	}
	defer autoStart.Close()
	reader := bufio.NewReader(autoStart)
	auto, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			autoStart.WriteString("false\n")
		}
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if strings.TrimSpace(auto) == "true" {
		return true
	}
	if strings.TrimSpace(auto) == "false" {
		return false
	}
	return false
}

//今日自动打卡是否成功,请确保auto_start.config文件存在
func TodayClockInSuccess() bool {
	clockInInfo, err := util.ReverseRead("./auto_start.config", 2)
	if err != nil {
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if len(clockInInfo) == 1 {
		if clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
			return true
		}
	}
	if strings.TrimSpace(clockInInfo[0]) == "" {
		return clockInInfo[1] == time.Now().Format("2006/01/02")+" 自动打卡成功"
	}
	if clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
		return true
	}
	return false
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

func WrittenToTheLog(content string) {
	CheckAutoStart()
	auto_start, err := os.OpenFile("./auto_start.config", os.O_APPEND|os.O_WRONLY, 0666)
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
