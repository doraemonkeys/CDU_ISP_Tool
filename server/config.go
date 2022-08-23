package server

import (
	"ISP_Tool/model"
	"ISP_Tool/utils"
	"ISP_Tool/view"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func InitConfig() error {
	log.Println("正在初始化配置文件")
	fmt.Println("正在初始化配置文件")
	err := os.MkdirAll("./config", 0666)
	if err != nil {
		log.Println("创建程序配置文件夹失败！", err)
		fmt.Println("创建程序配置文件夹失败！", err)
	}
	model.Auto_Start = CheckAutoStart()
	if model.Auto_Start {
		if TodayClockInSuccess() {
			//已经设置为自启动并且今日打卡已成功
			log.Println("打卡程序重新运行的原因: 用户手动打开")
			fmt.Println()
			model.Auto_Clock_IN_Success = true
			view.Auto_Clock_IN_Success()
			fmt.Println()
			startTime := time.Now()
			fmt.Println()
			fmt.Printf("按Enter键继续执行程序......")
			ch := make(chan bool, 1)
			go utils.PressToContinue(ch)
			ok := false
			for !ok {
				select {
				case ok = <-ch:
				default:
					time.Sleep(time.Second / 4)
				}
				//无操作30秒退出
				if time.Since(startTime) > time.Minute/2 {
					os.Exit(0)
				}
			}
		}
	}
	fmt.Println("正在检查网络环境...")
	for i := 0; !utils.NetWorkStatus(); i++ {
		time.Sleep(time.Second)
		//最多检查10次
		if i == 10 {
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
	config, err := os.Open("./config/配置文件.config")
	if err == nil {
		defer config.Close()
		temp := make([]byte, 20)
		n, err := config.Read(temp)
		if err != nil && err != io.EOF {
			log.Println("预读取配置文件失败，Error:", err)
			fmt.Println("预读取配置文件失败，Error:", err)
			return err
		}
		if n < 10 {
			log.Println("配置文件为空!")
			fmt.Println("配置文件为空!")
		} else {
			return nil
		}
	} else {
		log.Println("配置文件不存在或打开失败", err)
		fmt.Println("配置文件不存在或打开失败")
	}
	err = firstUse()
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println()
	model.UserConfigChanged = true
	return nil
}

func RebuitConfig(users []model.UserInfo) error {
	config, err := os.OpenFile("./config/配置文件.config", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("打开配置文件失败，Error:", err)
		fmt.Println("打开配置文件失败，Error:", err)
		return err
	}
	defer config.Close()
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
	return nil
}

func SetAutoStart() error {
	startPath := `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup\isp_auto_start.vbs`
	file, err := os.OpenFile(startPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Println("创建或打开文件失败!", err)
		return err
	}
	defer file.Close()
	path, err := utils.GetCurrentPath()
	if err != nil {
		log.Println("获取当前文件目录失败！", err)
		return err
	}
	lastindex := strings.LastIndex(path, "\\")
	path = path[:lastindex]
	path = strings.Replace(path, `\`, `\\`, -1)
	_, err = file.WriteString(utils.Utf8ToANSI(`Set objShell = CreateObject("WScript.Shell")` + "\n"))
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	_, err = file.WriteString(utils.Utf8ToANSI(`objShell.CurrentDirectory = "` + path + `\\config` + `"` + "\n"))
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	_, err = file.WriteString(utils.Utf8ToANSI(`objShell.Run "powershell /c ` + ".\\*.exe" + `"` + `,0`))
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	start_config, err := os.OpenFile("./config/auto_start.config", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println("创建或打开自启动配置文件失败!", err)
		return err
	}
	defer start_config.Close()
	n, err := start_config.WriteAt([]byte("true "), 0)
	if err != nil || n != 5 {
		log.Println("写入自启动配置文件失败！", err)
		return err
	}
	err = StartNewProgram()
	if err != nil {
		return err
	}
	return nil
}

//用户设置自启动后会关闭当前程序，开启一个守护进程，
//应当确保在设置自启动后调用。
func StartNewProgram() error {
	//延迟几秒打开一个新进程,不等cmd执行完毕就返回
	_, err := utils.Cmd_NoWait(`C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup`,
		[]string{"ping -n 5 127.1>nul", "&", ".\\isp_auto_start.vbs"})
	if err != nil {
		log.Println("打开新的程序失败！", err)
		return err
	}
	return nil
}

func CancelAutoStart() error {
	startPath := `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup\isp_auto_start.vbs`
	file, err := os.OpenFile(startPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Println("创建或打开文件失败!", err)
		return err
	}
	defer file.Close()
	start_config, err := os.OpenFile("./config/auto_start.config", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println("创建或打开自启动配置文件失败!", err)
		return err
	}
	defer start_config.Close()
	n, err := start_config.WriteAt([]byte("false"), 0)
	if err != nil || n != 5 {
		log.Println("写入自启动配置文件失败！", err)
		return err
	}
	return nil
}

//添加用户信息
func AddUser() error {
	config, err := os.OpenFile("./config/配置文件.config", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("打开配置文件失败，Error:", err)
		fmt.Println("打开配置文件失败，Error:", err)
		return err
	}
	defer config.Close()
	var users []model.UserInfo
	var NewUser = model.UserInfo{}
	for {
		fmt.Println()
		attributes := [5]color.Attribute{}
		attributes[1] = color.FgRed
		utils.ColorPrint(attributes[:], "输入 ", "Q", " 退出添加账号\n")
		//fmt.Println("输入 Q 退出添加账号")
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