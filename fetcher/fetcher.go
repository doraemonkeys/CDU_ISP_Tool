package fetcher

import (
	"ISP_Tool/model"
	"ISP_Tool/util"
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

func Get_User_Nonce(client *http.Client) (string, error) {
	request, _ := http.NewRequest("GET", "https://xsswzx.cdu.edu.cn/ispstu/com_user/webindex.asp", nil)
	request.Header.Set("authority", "xsswzx.cdu.edu.cn")
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("user-agent", model.UserAgent)
	request.Header.Set("referer", "https://xsswzx.cdu.edu.cn/ispstu/com_user/weblogin.asp")

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
	re := regexp.MustCompile(model.User_no_RE)
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
	client, err := util.Get_client()
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
	ip, err := util.GetIPV4()
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
	e, err := util.DetermineEncodingbyPeek(bodyReader)
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
	re := regexp.MustCompile(model.Ip_locationRe)
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
	apiUrl := "https://xsswzx.cdu.edu.cn/ispstu/com_user/projecthealth.asp"
	// URL param
	queryData := url.Values{}
	queryData.Set("id", user_no)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return model.Location{}, errors.New("parse url requestUrl failed")
	}
	u.RawQuery = queryData.Encode() // URL encode

	request, _ := http.NewRequest("GET", u.String(), nil)
	request.Header.Set("authority", "xsswzx.cdu.edu.cn")
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("user-agent", model.UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		log.Println("访问ISP打开页面失败，可能是ISP结构发生变化，请联系开发者。返回的状态：", http.StatusText(resp.StatusCode))
		fmt.Println("访问ISP打开页面失败，可能是ISP结构发生变化，请联系开发者。返回的状态：", http.StatusText(resp.StatusCode))
		return model.Location{}, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("访问ISP打开页面失败！", err)
		fmt.Println("访问ISP打开页面失败！", err)
		return model.Location{}, err
	}
	re := regexp.MustCompile(model.Isp_history_location_Re)
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

func GetLocation(user_no string, client *http.Client) (model.Location, error) {
	isp_location_history, err1 := Get_isp_location_history(user_no, client)
	if err1 != nil {
		log.Println("获取isp历史打卡信息失败！")
		fmt.Println("获取isp历史打卡信息失败！")
	} else {
		fmt.Println("历史健康登记打卡地址：",
			isp_location_history.Province, isp_location_history.City, isp_location_history.Area)
		log.Println("历史健康登记打卡地址：",
			isp_location_history.Province, isp_location_history.City, isp_location_history.Area)
	}
	IP_Loaction, err2 := Get_IP_Loaction()
	if err2 != nil {
		log.Println("获取ip地址信息失败！")
		fmt.Println("获取ip地址信息失败！")
	} else {
		fmt.Printf("当前ip地址：")
		color.Yellow("%s %s %s", IP_Loaction.Province, IP_Loaction.City, IP_Loaction.Area)
		log.Println("当前ip地址：",
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
