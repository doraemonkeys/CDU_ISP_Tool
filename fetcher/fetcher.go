package fetcher

import (
	"ISP_Tool/model"
	"ISP_Tool/server"
	"ISP_Tool/utils"
	"bufio"
	"errors"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/fatih/color"
	"golang.org/x/text/transform"
)

//检查打卡是否异常
func CheckingAnomalies(user_no string, client *http.Client) (model.FieldAndValue, error) {
	apiUrl := model.All.ClockInHome.ClockInHomeUrl
	// URL param
	queryData := url.Values{}
	queryData.Set(model.All.ClockInHome.QueryField, user_no)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return model.FieldAndValue{}, errors.New("parse url requestUrl failed")
	}
	u.RawQuery = queryData.Encode() // URL encode

	request, _ := http.NewRequest(model.All.ClockInHome.Head.Method, u.String(), nil)
	request.Header.Set("authority", model.All.ClockInHome.Head.Authority)
	request.Header.Set("content-type", model.All.ClockInHome.Head.Content_type)
	request.Header.Set("user-agent", model.UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		log.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		return model.FieldAndValue{}, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("访问ISP页面失败！", err)
		fmt.Println("访问ISP页面失败！", err)
		return model.FieldAndValue{}, err
	}
	//匹配撤回打卡key-value字段
	var key_value = model.FieldAndValue{}
	re2 := regexp.MustCompile(model.All.Regexp.Clock_IN_ID_Re)
	match2 := re2.FindSubmatch(content)
	if len(match2) != 3 {
		log.Println("获取撤回打卡key-value字段失败,可能是ISP结构发生变化，请联系开发者。")
	} else {
		key_value.Field = string(match2[1])
		key_value.Value = string(match2[2])
		log.Printf("获取到当前打卡%s(ID)：%s\n", key_value.Field, key_value.Value)
	}
	err = server.LookForKeyword(content)
	if err != nil {
		log.Println("健康登记出现异常，可能是程序的错误或正处于风险区！")
		color.Red("健康登记出现异常，可能是程序的错误或正处于风险区！")
		return key_value, err
	}
	return key_value, nil
}

func Get_User_Nonce(client *http.Client) (string, error) {
	request, _ := http.NewRequest(model.All.ISPHome.Head.Method, model.All.ISPHome.ISPHomeUrl, nil)
	request.Header.Set("authority", model.All.ISPHome.Head.Authority)
	request.Header.Set("content-type", model.All.ISPHome.Head.Content_type)
	request.Header.Set("user-agent", model.UserAgent)
	request.Header.Set("referer", model.All.ISPHome.Head.Referer)

	resp, err := client.Do(request)
	if err != nil {
		log.Println("访问ISP主页(webindex.asp)失败，可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("访问ISP主页(webindex.asp)失败，可能是ISP结构发生变化，请联系开发者。")
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取ISP主页(webindex.asp)失败！", err)
		fmt.Println("读取ISP主页(webindex.asp)失败！", err)
		return "", err
	}
	re := regexp.MustCompile(model.All.Regexp.User_no_Re)
	match := re.FindSubmatch(content)
	var user_no string
	if match == nil {
		log.Println("访问ISP主页寻找用户标识码失败，可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("访问ISP主页寻找用户标识码失败，可能是ISP结构发生变化，请联系开发者。")
		return "", errors.New("match == nil")
	}
	user_no = string(match[1])
	return user_no, nil
}

func Get_IP_Loaction() (model.Location, error) {
	client, err := utils.Get_client()
	if err != nil {
		log.Println("程序初始化client失败，请联系开发者。", err)
		fmt.Println("程序初始化client失败，请联系开发者。", err)
		return model.Location{}, err
	}
	apiUrl := "https://www.ip138.com/iplookup.asp"
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return model.Location{}, err
	}
	// URL param
	queryData := url.Values{}
	ip, err := utils.GetIPV4()
	if err != nil {
		return model.Location{}, err
	}
	fmt.Println("获取到公网IPv4:", ip)
	log.Println("获取到公网IPv4:", ip)
	queryData.Set("ip", ip)
	queryData.Set("action", "2")
	u.RawQuery = queryData.Encode() // URL encode
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("User-Agent", model.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("访问ip138.com失败！, err:%v\n", err)
		log.Printf("访问ip138.com失败！, err:%v\n", err)
		return model.Location{}, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	//自动检测html编码
	e, err := utils.DetermineEncodingbyPeek(bodyReader)
	if err != nil {
		log.Println("检测html编失败，请联系开发者。", err)
		fmt.Println("检测html编失败，请联系开发者。", err)
		return model.Location{}, err
	}
	//转码utf-8
	utf8BodyReader := transform.NewReader(bodyReader, e.NewDecoder())

	content, err := ioutil.ReadAll(utf8BodyReader)
	if err != nil {
		log.Println("读取网页ip138.com失败！", err)
		fmt.Println("读取网页ip138.com失败！", err)
		return model.Location{}, err
	}
	re := regexp.MustCompile(model.All.Regexp.Ip_locationRe)
	match := re.FindSubmatch(content)
	if match == nil {
		log.Println("匹配IP地址信息失败！请联系开发者。")
		fmt.Println("匹配IP地址信息失败！请联系开发者。")
		return model.Location{}, errors.New("match == nil")
	}
	if len(match) != 5 {
		log.Println("匹配IP地址信息失败！请联系开发者。")
		fmt.Println("匹配IP地址信息失败！请联系开发者。")
		return model.Location{}, errors.New("match == nil")
	}
	newLocation := model.Location{}
	if string(match[1]) != "中国" {
		log.Println("匹配IP地址信息失败！", "国家：", string(match[1]))
		fmt.Println("匹配IP地址信息失败！", "国家：", string(match[1]))
		return model.Location{}, errors.New("match[1] != 中国")
	}
	newLocation.Province = string(match[2])
	newLocation.City = string(match[3])
	newLocation.Area = string(match[4])
	return newLocation, nil
}

func Get_isp_location_history(user_no string, client *http.Client) (model.Location, error) {
	apiUrl := model.All.ClockInHome.ClockInHomeUrl
	// URL param
	queryData := url.Values{}
	queryData.Set(model.All.ClockInHome.QueryField, user_no)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return model.Location{}, errors.New("parse url requestUrl failed")
	}
	u.RawQuery = queryData.Encode() // URL encode

	request, _ := http.NewRequest(model.All.ClockInHome.Head.Method, u.String(), nil)
	request.Header.Set("authority", model.All.ClockInHome.Head.Authority)
	request.Header.Set("content-type", model.All.ClockInHome.Head.Content_type)
	request.Header.Set("user-agent", model.UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		log.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		return model.Location{}, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取ISP页面内容失败！", err)
		fmt.Println("读取ISP页面内容失败！", err)
		return model.Location{}, err
	}
	re := regexp.MustCompile(model.All.Regexp.Isp_history_location_Re)
	match := re.FindSubmatch(content)
	if match == nil {
		log.Println("匹配ISP历史登记地址信息失败！可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("匹配ISP历史登记地址信息失败！可能是ISP结构发生变化，请联系开发者。")
		return model.Location{}, err
	}
	newLocation := model.Location{}
	newLocation.Province = string(match[1])
	newLocation.City = string(match[2])
	newLocation.Area = string(match[3])
	return newLocation, nil
}

func GetLocation(user model.UserInfo, client *http.Client) (model.Location, error) {
	isp_location_history, err1 := Get_isp_location_history(user.UserNonce, client)
	if err1 != nil {
		log.Println("获取isp历史打卡信息失败！")
		fmt.Println("获取isp历史打卡信息失败！")
		user.ChooseLocation = 1 //只能选择使用ip地址了
	} else {
		fmt.Printf("历史健康登记打卡地址：")
		color.Yellow("%s %s %s", isp_location_history.Province, isp_location_history.City, isp_location_history.Area)
		log.Println("历史健康登记打卡地址：",
			isp_location_history.Province, isp_location_history.City, isp_location_history.Area)
	}
	IP_Loaction, err2 := Get_IP_Loaction()
	if err2 != nil {
		log.Println("获取ip地址信息失败！")
		fmt.Println("获取ip地址信息失败！")
		user.ChooseLocation = 2 //只能选择使用isp历史打卡地址了
	} else {
		fmt.Printf("当前ip地址：     ")
		color.Yellow("%s %s %s", IP_Loaction.Province, IP_Loaction.City, IP_Loaction.Area)
		log.Println("当前ip地址：",
			IP_Loaction.Province, IP_Loaction.City, IP_Loaction.Area)
	}
	//全部出错
	if err1 != nil && err2 != nil {
		log.Println("获取地址信息失败,无法打卡！")
		fmt.Println("获取地址信息失败,无法打卡！")
		return model.Location{}, errors.New(err1.Error() + err2.Error())
	}
	attributes := [5]color.Attribute{}
	attributes[0] = color.FgYellow
	//全部没错或有一个获取失败,获取失败会导致user.ChooseLocation被修改(仅此函数中)
	if err2 == nil && user.ChooseLocation == 1 {
		utils.ColorPrint(attributes[:], "使用ip地址信息打卡", "，如果有错误请前往ISP手动修改！\n")
		log.Println("使用ip地址信息打卡，如果有错误请前往ISP手动修改！")
		return IP_Loaction, nil
	} else {
		utils.ColorPrint(attributes[:], "使用ISP历史登记地址打卡", "，如果有错误请前往ISP手动修改！\n")
		log.Println("使用ISP历史登记地址打卡，如果有错误请前往ISP手动修改！")
		return isp_location_history, nil
	}
}
