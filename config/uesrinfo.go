package config

import (
	"ISP_Tool/model"
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
	for {
		userData, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(userData) > 1 {
				userData = strings.TrimSpace(userData)
				json.Unmarshal([]byte(userData), &user)
				if user.UserID == NewUser.UserID {
					user.UserPwd = NewUser.UserPwd
					users = append(users, user)
				} else {
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
		if user.UserID == NewUser.UserID {
			user.UserPwd = NewUser.UserPwd
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
