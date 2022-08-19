package config

import (
	"ISP_Tool/model"
	"ISP_Tool/util"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func GetUserInfos() ([]model.UserInfo, error) {
	config, err := os.Open("./配置文件.config")
	if err != nil {
		log.Println("打开配置文件失败！", err)
		fmt.Println("打开配置文件失败！", err)
		return []model.UserInfo{}, err
	}
	defer config.Close()
	reader := bufio.NewReader(config)
	users := []model.UserInfo{}
	var user model.UserInfo
	for {
		userData, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(userData) > 1 {
				userData = strings.TrimSpace(userData)
				json.Unmarshal([]byte(userData), &user)
				users = append(users, user)
			}
			break
		}
		if err != nil {
			log.Println("读取配置文件失败！", err)
			fmt.Println("读取配置文件失败！", err)
			return []model.UserInfo{}, err
		}
		userData = strings.TrimSpace(userData)
		json.Unmarshal([]byte(userData), &user)
		users = append(users, user)
	}
	return users, nil
}

func ModifyUserInfos() error {
	var NewUser model.UserInfo
	fmt.Println("输入 Q 退出修改账号")
	fmt.Println("请输入要修改的学号：")
	var id string
	fmt.Scan(&id)
	NewUser.UserID = strings.TrimSpace(id)
	if NewUser.UserID == "Q" || NewUser.UserID == "q" {
		return errors.New("取消修改")
	}
	fmt.Println("请输入新的密码：")
	var pwd string
	fmt.Scan(&pwd)
	NewUser.UserPwd = strings.TrimSpace(pwd)
	if NewUser.UserPwd == "Q" || NewUser.UserPwd == "q" {
		return errors.New("取消修改")
	}
	config, err := os.Open("./配置文件.config")
	if err != nil {
		log.Println("打开配置文件失败！", err)
		fmt.Println("打开配置文件失败！", err)
		return err
	}
	defer config.Close()
	reader := bufio.NewReader(config)
	users := []model.UserInfo{}
	var user model.UserInfo
	found := false
	for {
		userData, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(userData) > 1 {
				userData = strings.TrimSpace(userData)
				json.Unmarshal([]byte(userData), &user)
				if user.UserID == NewUser.UserID {
					user.UserPwd = NewUser.UserPwd
					found = true
				}
				users = append(users, user)
			}
			if !found {
				return errors.New("没有找到目标ID")
			}
			break
		}
		if err != nil {
			log.Println("读取配置文件失败！", err)
			fmt.Println("读取配置文件失败！", err)
			return err
		}
		userData = strings.TrimSpace(userData)
		json.Unmarshal([]byte(userData), &user)
		if user.UserID == NewUser.UserID {
			user.UserPwd = NewUser.UserPwd
			found = true
		}
		users = append(users, user)
	}
	err = RebuitConfig(users)
	if err != nil {
		log.Println("修改配置文件失败！", err)
		fmt.Println("修改配置文件失败！", err)
		return err
	}
	return nil
}

func DeleteUser() error {
	var targetUser model.UserInfo
	fmt.Println("输入 Q 退出修改账号")
	fmt.Println("请输入要删除的学号：")
	var id string
	fmt.Scan(&id)
	targetUser.UserID = strings.TrimSpace(id)
	if targetUser.UserID == "Q" || targetUser.UserID == "q" {
		return errors.New("取消删除")
	}
	config, err := os.Open("./配置文件.config")
	if err != nil {
		log.Println("打开配置文件失败！", err)
		fmt.Println("打开配置文件失败！", err)
		return err
	}
	defer config.Close()
	reader := bufio.NewReader(config)
	users := []model.UserInfo{}
	var user model.UserInfo
	for {
		userData, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(userData) > 1 {
				userData = strings.TrimSpace(userData)
				json.Unmarshal([]byte(userData), &user)
				if user.UserID != targetUser.UserID {
					return errors.New("没有找到目标ID")
				}
			}
			break
		}
		if err != nil {
			log.Println("读取配置文件失败！", err)
			fmt.Println("读取配置文件失败！", err)
			return err
		}
		userData = strings.TrimSpace(userData)
		json.Unmarshal([]byte(userData), &user)
		if user.UserID == targetUser.UserID {
			continue
		}
		users = append(users, user)
	}
	err = RebuitConfig(users)
	if err != nil {
		log.Println("修改配置文件失败！", err)
		fmt.Println("修改配置文件失败！", err)
		return err
	}
	return nil
}

func RebuitConfig(users []model.UserInfo) error {
	config, err := os.OpenFile("./配置文件.config", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
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
	path, err := util.GetCurrentPath()
	if err != nil {
		log.Println("获取当前文件目录失败！", err)
		return err
	}
	lastindex := strings.LastIndex(path, "\\")
	toolName := path[lastindex+1:]
	path = path[:lastindex]
	path = strings.Replace(path, `\`, `\\`, -1)
	_, err = file.WriteString(`Set objShell = CreateObject("WScript.Shell")` + "\n")
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	_, err = file.WriteString(`objShell.CurrentDirectory = "` + path + `"` + "\n")
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	_, err = file.WriteString(`objShell.Run "cmd /c ` + toolName + `"` + `,0`)
	if err != nil {
		log.Println("写入当前文件目录失败！", err)
		return err
	}
	start_config, err := os.OpenFile("./auto_start.config", os.O_CREATE|os.O_RDWR, 0666)
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

//用户设置自启动后会关闭当前程序，所以要开一个新的进程，
//应当确保在设置自启动后调用。
func StartNewProgram() error {
	//延迟7秒打开一个新进程,不等cmd执行完毕就返回
	_, err := util.Cmd_NoWait(`C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup`,
		[]string{"ping -n 7 127.1>nul", "&", ".\\isp_auto_start.vbs"})
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
	start_config, err := os.OpenFile("./auto_start.config", os.O_CREATE|os.O_RDWR, 0666)
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
		fmt.Println()
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
