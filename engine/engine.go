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
)

func Run(client http.Client, user model.UserInfo) error {
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
	userLocation, err := getLocation(user_no, &client)
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

func getLocation(user_no string, client *http.Client) (model.Location, error) {
	isp_location_history, err1 := fetcher.Get_isp_location_history(user_no, client)
	if err1 != nil {
		log.Println("获取isp历史打卡信息失败！")
		fmt.Println("获取isp历史打卡信息失败！")
	} else {
		fmt.Println("历史健康登记打卡地址：",
			isp_location_history.Province, isp_location_history.City, isp_location_history.Area)
	}
	IP_Loaction, err2 := fetcher.Get_IP_Loaction()
	if err2 != nil {
		log.Println("获取ip地址信息失败！")
		fmt.Println("获取ip地址信息失败！")
	} else {
		fmt.Println("当前ip地址：",
			IP_Loaction.Province, IP_Loaction.City, IP_Loaction.Area)
	}
	if err1 != nil && err2 != nil {
		log.Println("获取地址信息失败,无法打开！")
		fmt.Println("获取地址信息失败,无法打开！")
		return model.Location{}, errors.New(err1.Error() + err2.Error())
	}
	if err2 == nil {
		fmt.Println("默认使用ip地址信息打卡，如果有错误请前往ISP手动修改！")
		log.Println("默认使用ip地址信息打卡，如果有错误请前往ISP手动修改！")
		return IP_Loaction, nil
	}
	fmt.Println("使用ISP历史登记信息打卡，如果有错误请前往ISP手动修改！")
	log.Println("使用ISP历史登记信息打卡，如果有错误请前往ISP手动修改！")
	return isp_location_history, nil
}
