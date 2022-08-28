package server

import (
	"ISP_Tool/model"
	"ISP_Tool/utils"
	"ISP_Tool/view"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

//用户可能是第一次使用
func firstUse() error {
	config, err := os.OpenFile("./config/配置文件.config", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
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
	NewUser.ChooseLocation = 1
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
	return nil
}

func ReachConsensus() bool {
	fmt.Println(">>>>>>>>>>>>没有检测到账号，你可能是第一次使用本软件！<<<<<<<<<<<<<<<<")
	fmt.Println(">>>>>>>>>>>>          请阅读注意事项                  <<<<<<<<<<<<<<<<")
	fmt.Println()
	view.Warn()
	fmt.Println()
	var input string
	//fmt.Println("是否同意？输入 [" + color.RedString("Yes") + "] or [" + color.RedString("No") + "]")
	attributes := [5]color.Attribute{}
	attributes[1] = color.FgRed
	attributes[3] = color.FgRed
	utils.ColorPrint(attributes[:], "是否同意？输入 [", "Yes", "] or [", "No", "]\n")
	fmt.Scan(&input)
	input = strings.TrimSpace(input)
	input = strings.ToUpper(input)
	for input != "YES" {
		if input == "NO" {
			return false
		}
		utils.ColorPrint(attributes[:], "是否同意？输入 [", "Yes", "] or [", "No", "]\n")
		fmt.Scan(&input)
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)
	}
	return true
}

func GetUserInfos() ([]model.UserInfo, error) {
	config, err := os.Open("./config/配置文件.config")
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
	config, err := os.Open("./config/配置文件.config")
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
	config, err := os.Open("./config/配置文件.config")
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
				if user.UserID == targetUser.UserID {
					found = true
				} else {
					users = append(users, user)
				}
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
		if user.UserID == targetUser.UserID {
			found = true
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
