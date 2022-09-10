package uploader

import (
	"ISP_Tool/model"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ISP_CheckIn(client *http.Client, user model.UserInfo) error {
	today := time.Now().Local().Format("2006年1月2日")
	apiUrl := model.All.ClockIn.ClockInUrl
	//URL param
	queryData := url.Values{}
	queryData.Set("id", user.UserNonce)
	queryData.Set("id2", today)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return err
	}
	u.RawQuery = queryData.Encode() // URL encode

	// 构造请求
	parames, err := getPostField(user, client)
	if err != nil {
		return err
	}
	data := parames.Encode()
	req4, _ := http.NewRequest(model.All.ClockIn.Head.Method, u.String(), strings.NewReader(data))
	req4.Header.Set("authority", model.All.ClockIn.Head.Authority)
	req4.Header.Set("content-type", model.All.ClockIn.Head.Content_type)
	req4.Header.Set("referer", model.All.ClockIn.Head.Referer)
	req4.Header.Set("user-agent", model.UserAgent)

	resp4, err := client.Do(req4)
	if err != nil {
		log.Println("发起ISP登记请求失败！可能是ISP结构发生变化，请联系开发者。", err)
		fmt.Println("发起ISP登记请求失败！可能是ISP结构发生变化，请联系开发者。", err)
		return err
	}
	content4, err := ioutil.ReadAll(resp4.Body)
	if err != nil {
		log.Println("读取ISP登记返回信息失败！", err)
		fmt.Println("读取ISP登记返回信息失败！", err)
		return err
	}
	//fmt.Println(string(content4))
	re := regexp.MustCompile(model.All.Regexp.Clock_In_success_Re)
	match := re.Find(content4)
	if match != nil {
		return nil
	}
	re = regexp.MustCompile(model.All.Regexp.Already_Clock_In_Re)
	match = re.Find(content4)
	if match != nil {
		log.Println(user.UserID, "健康登记打卡已存在")
		color.Green("%s %s", user.UserID, "健康登记打卡已存在")
		return errors.New("健康登记打卡已存在")
	}
	return errors.New("CDU-ISP 健康登记打卡 失败")
}

func getPostField(user model.UserInfo, client *http.Client) (url.Values, error) {
	today := time.Now().Local().Format("2006年1月2日")
	apiUrl := model.All.ClockIn.ClockInUrl
	//URL param
	queryData := url.Values{}
	queryData.Set("id", user.UserNonce)
	queryData.Set("id2", today)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return nil, err
	}
	u.RawQuery = queryData.Encode() // URL encode

	// 构造请求
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("authority", model.All.ClockIn.Head.Authority)
	req.Header.Set("content-type", model.All.ClockIn.Head.Content_type)
	req.Header.Set("referer", model.All.ClockIn.Head.Referer)
	req.Header.Set("user-agent", model.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("访问ISP登记请求页面失败！可能是ISP结构发生变化，请联系开发者。", err)
		fmt.Println("访问ISP登记请求页面失败！可能是ISP结构发生变化，请联系开发者。", err)
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取ISP登记请求页面失败！", err)
		fmt.Println("读取ISP登记请求页面失败！", err)
		return nil, err
	}
	var parame = url.Values{}
	re := regexp.MustCompile(model.All.ClockIn.MatchActionRe)
	match := re.FindSubmatch(content)
	if match == nil {
		log.Println("匹配登记请求字段失败！")
		fmt.Println("匹配登记请求字段失败！")
		return nil, errors.New("匹配登记请求字段失败")
	}
	parame.Add("action", string(match[1]))
	re = regexp.MustCompile(model.All.ClockIn.MatchParamesRe)
	match2 := re.FindAllSubmatch(content, -1)
	if match2 == nil {
		log.Println("匹配登记请求字段失败！")
		fmt.Println("匹配登记请求字段失败！")
		return nil, errors.New("匹配登记请求字段失败")
	}
	//先全部否定
	for _, v := range match2 {
		parame.Add(string(v[1]), "否")
	}
	//设置区域
	parame.Set(model.All.ClockIn.AreaField, user.Area)
	parame.Set(model.All.ClockIn.CityField, user.City)
	parame.Set(model.All.ClockIn.ProvinceField, user.Province)
	//打印parame的所有key
	count := 0
	for k, v := range parame {
		for range v {
			fmt.Printf("%s ", k)
			count++
		}
	}
	//打印parame的数量
	fmt.Println("共", count, "个字段")
	return parame, nil
}

func CancelCheckIn(key_value model.FieldAndValue, client *http.Client) error {
	success := false
	var err error
	for i := 0; !success; i++ {
		err = tryCancle(key_value, client)
		if err == nil {
			return nil
		}
		//只尝试取消10次
		if i > 10 {
			break
		}
		time.Sleep(time.Second / 4)
	}
	return err
}

func tryCancle(key_value model.FieldAndValue, client *http.Client) error {
	//apiUrl := "https://xsswzx.cdu.edu.cn/ispstu/com_user/projecthealth_del.asp"
	apiUrl := model.All.Cancel.CancelUrl
	// URL param
	queryData := url.Values{}
	queryData.Set(key_value.Field, key_value.Value)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return errors.New("parse url requestUrl failed")
	}
	u.RawQuery = queryData.Encode() // URL encode

	request, _ := http.NewRequest(model.All.Cancel.Head.Method, u.String(), nil)
	//request.Header.Set("authority", "xsswzx.cdu.edu.cn")
	//request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("user-agent", model.UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		log.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("访问ISP页面失败，可能是ISP结构发生变化，请联系开发者。")
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取ISP页面内容失败！", err)
		fmt.Println("读取ISP页面内容失败！", err)
		return err
	}
	re := regexp.MustCompile(model.All.Regexp.CancelSuccessRe)
	match := re.Find(content)
	if match != nil {
		return nil
	}
	return errors.New("发送删除请求成功,但出现未知的错误！")
}
