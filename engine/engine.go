package engine

import (
	"ISP_Tool/fetcher"
	"ISP_Tool/login"
	"ISP_Tool/model"
	"ISP_Tool/uploader"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
)

func Run(client http.Client, user model.UserInfo) error {
	if !model.UserConfigChanged && model.Auto_Start && model.Auto_Clock_IN_Success {
		log.Println(user.UserID, "健康登记打卡已存在")
		color.Green("%s %s", user.UserID, "健康登记打卡已存在")
		return errors.New("健康登记打卡已存在")
	}
	err := login.LoginISP(&client, user)
	if err != nil {
		fmt.Println("ID", user.UserID, "登录失败！")
		return err
	}
	log.Println("登录ISP成功")
	fmt.Println("登录ISP成功")
	user_no, err := fetcher.Get_User_Nonce(&client)
	if err != nil {
		fmt.Println("ID", user.UserID, "获取ISP的用户识别码失败！")
		return err
	}
	fmt.Println("从isp获取到用户识别码：", user_no)
	log.Println("从isp获取到用户识别码：", user_no)
	user.UserNonce = user_no
	userLocation, err := fetcher.GetLocation(user_no, &client)
	if err != nil {
		fmt.Println("ID", user.UserID, "获取地理位置失败！")
		return err
	}
	user.Location = userLocation
	err = uploader.ISP_Clock_In(&client, user)
	if err != nil {
		return err
	}
	return nil
}
