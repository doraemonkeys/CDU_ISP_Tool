package login

import (
	"ISP_Tool/model"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dop251/goja"
	"github.com/fatih/color"
)

// 获取预登录信息
func getLoginInfo(client *http.Client) (string, string, string, error) {
	request, err := http.NewRequest("GET", "https://vpn.cdu.edu.cn/por/login_auth.csp?apiversion=1", nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return "", "", "", err
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", "err:", err)
		return "", "", "", err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return "", "", "", err
	}
	re := regexp.MustCompile(`RAND_CODE>([^<]{2,}?)<`)
	match := re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配RAND_CODE失败")
		return "", "", "", errors.New("匹配RAND_CODE失败")
	}
	CSRF_RAND_CODE := string(match[1])
	//fmt.Println(CSRF_RAND_CODE)
	re = regexp.MustCompile(`RSA_ENCRYPT_KEY>([^<]{2,}?)<`)
	match = re.FindSubmatch(content)

	if len(match) == 0 {
		log.Println("匹配RSA_ENCRYPT_KEY失败")
		return "", "", "", errors.New("匹配RSA_ENCRYPT_KEY失败")
	}
	RSA_ENCRYPT_KEY := string(match[1])
	//fmt.Println(RSA_ENCRYPT_KEY)
	re = regexp.MustCompile(`RSA_ENCRYPT_EXP>([^<]{2,}?)<`)
	match = re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配RSA_ENCRYPT_EXP失败")
		return "", "", "", errors.New("匹配RSA_ENCRYPT_EXP失败")
	}
	RSA_ENCRYPT_EXP := string(match[1])
	//fmt.Println(RSA_ENCRYPT_EXP)
	// re = regexp.MustCompile(`TwfID>([^<]{2,}?)<`)
	// match = re.FindSubmatch(content)
	// if len(match) == 0 {
	// 	log.Println("匹配TwfID失败")
	// 	return nil, 0, errors.New("匹配TwfID失败")
	// }
	// TwfID := string(match[1])
	// fmt.Println(TwfID)
	return CSRF_RAND_CODE, RSA_ENCRYPT_KEY, RSA_ENCRYPT_EXP, nil
}

func getField_V(client *http.Client) (string, error) {
	//访问https://vpn.cdu.edu.cn/portal/
	request, err := http.NewRequest("GET", "https://vpn.cdu.edu.cn/portal/", nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return "", err
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", "err:", err)
		return "", err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return "", err
	}
	re := regexp.MustCompile(`v[\s]*=[\s]*([0-9]+)`)
	match := re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配v失败")
		return "", errors.New("匹配v失败")
	}
	v := string(match[1])
	//fmt.Println(v)
	return v, nil
}

// 导出JS函数
func exporetJSFunc(client *http.Client, v string) (func(string, string, string) string, error) {
	jsPrefix := `
	var navigator = navigator || {};
	var window = window || {};
	`
	jsSuffix := `
	function rasPwd(r, i, pwd) {
		var t, n, r1, i1, o, s = new RSAKey;
		i1 = i;
		r1 = r;
		s.setPublic(r1, i1);
		t = s.encrypt(pwd);
		return t;
	}
	`
	//get https://vpn.cdu.edu.cn/portal/libs/rsa.js?v=582537319
	request, err := http.NewRequest("GET", "https://vpn.cdu.edu.cn/portal/libs/rsa.js?v="+v, nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", err)
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return nil, err
	}
	jsFile := jsPrefix + string(content) + jsSuffix
	vm := goja.New()
	_, err = vm.RunString(jsFile)
	if err != nil {
		log.Println("执行js失败", err)
		return nil, err
	}
	//r i pwd,r就是RSA_ENCRYPT_KEY,i为定值,pwd为密码+"_"+randcode
	var fn func(string, string, string) string
	err = vm.ExportTo(vm.Get("rasPwd"), &fn)
	if err != nil {
		log.Println("导出函数失败", err)
		return nil, err
	}
	return fn, nil
}

// Login VPN
func loginVPN(client *http.Client, postUrl string, postData url.Values) error {
	//POST
	request, err := http.NewRequest("POST", postUrl, strings.NewReader(postData.Encode()))
	if err != nil {
		log.Println("创建请求失败", err)
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", "err:", err)
		return err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return err
	}
	//fmt.Println(string(content))
	re := regexp.MustCompile(`password`)
	match := re.Find(content)
	if len(match) != 0 {
		log.Println("登录VPN失败,账号或密码错误")
		return errors.New("登录VPN失败,账号或密码错误")
	}
	return nil
}

// Get VPN Login Post URL and postData
func getPost_URL_postData(client *http.Client, v string) (string, url.Values, error) {
	//GET https://vpn.cdu.edu.cn/portal/jssdk/api/auth_psw.js
	//prrameter: v=15975369
	request, err := http.NewRequest("GET", "https://vpn.cdu.edu.cn/portal/jssdk/api/auth_psw.js", nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return "", nil, err
	}
	q := request.URL.Query()
	q.Add("v", v)
	request.URL.RawQuery = q.Encode()
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", err)
		return "", nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return "", nil, err
	}
	re := regexp.MustCompile(`path[\s]*:[\s]*"([^"]+)"[\s\S]{1,50}post`)
	match := re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配path失败")
		return "", nil, errors.New("匹配path失败")
	}
	path := string(match[1])
	//fmt.Println(path)
	postUrl := "https://vpn.cdu.edu.cn" + path
	//fmt.Println(postUrl)
	re = regexp.MustCompile(`authLoginPwd[\s]*\(([^)]{10,}?)\)`)
	match = re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配authLoginPwd失败")
		return "", nil, errors.New("匹配authLoginPwd失败")
	}
	subContent := string(match[1])
	//fmt.Println(subContent)
	re = regexp.MustCompile(`([\w]+)[\s]*:[\s]*('|"|)([\w]*)('|"|)`)
	match2 := re.FindAllSubmatch([]byte(subContent), -1)
	if len(match2) == 0 {
		log.Println("匹配authLoginPwd中的字段失败")
		return "", nil, errors.New("匹配authLoginPwd中的字段失败")
	}
	postData := url.Values{}
	for _, v := range match2 {
		if string(v[2]) != "" {
			postData.Set(string(v[1]), string(v[3]))
		} else {
			postData.Set(string(v[1]), "")
		}
	}
	//fmt.Println(postData)
	return postUrl, postData, nil
}

func getISP_Port(client *http.Client) (string, int, error) {
	//get
	request, err := http.NewRequest("GET", "https://vpn.cdu.edu.cn/por/conf.csp?apiversion=1", nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return "", 0, err
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", err)
		return "", resp.StatusCode, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return "", resp.StatusCode, err
	}
	re := regexp.MustCompile(`WebVpn[\s\S]{1,60}?port[\s]*=['"]([0-9]+)['"]`)
	match := re.FindSubmatch(content)
	if len(match) == 0 {
		log.Println("匹配port失败")
		return "", resp.StatusCode, errors.New("匹配port失败")
	}
	port := string(match[1])
	//fmt.Println(port)
	model.All.VPN_ISP_BaseURL = model.All.VPN_ISP_BaseURL + ":" + port
	return port, resp.StatusCode, nil
}

func fetchISPByVPN(client *http.Client, port string) ([]byte, int, error) {
	//get http://xsswzx-cdu-edu-cn-s.vpn.cdu.edu.cn:8118/ispstu/com_user/weblogin.asp
	request, err := http.NewRequest("GET", "http://xsswzx-cdu-edu-cn-s.vpn.cdu.edu.cn:"+port+"/ispstu/com_user/weblogin.asp", nil)
	if err != nil {
		log.Println("创建请求失败", err)
		return nil, 0, err
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("请求失败", err)
		return nil, resp.StatusCode, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败", err)
		return nil, resp.StatusCode, err
	}
	return content, resp.StatusCode, nil
}

func Fetch_ISP_Login_Pag_VPN(client *http.Client, user model.UserInfo) ([]byte, int, error) {
	if user.VPN_Pwd == "" {
		color.Yellow("检测到VPN密码为空,将尝试使用ISP密码登录")
		user.VPN_Pwd = user.UserPwd
	}
	CSRF_RAND_CODE, RSA_ENCRYPT_KEY, _, err := getLoginInfo(client)
	if err != nil {
		log.Println("获取预登录信息失败", err)
		return nil, 0, err
	}
	v, err := getField_V(client)
	if err != nil {
		return nil, 0, err
	}
	fn, err := exporetJSFunc(client, v)
	if err != nil {
		return nil, 0, err
	}
	//pwd+rand
	pwd := user.VPN_Pwd + "_" + CSRF_RAND_CODE
	svpn_password := fn(RSA_ENCRYPT_KEY, "10001", pwd)
	//fmt.Println(svpn_password)
	postUrl, postData, err := getPost_URL_postData(client, v)
	if err != nil {
		return nil, 0, err
	}
	//遍历postData
	for k := range postData {
		if strings.Contains(k, "randcode") {
			postData.Set(k, CSRF_RAND_CODE)
			continue
		}
		if strings.Contains(k, "name") {
			postData.Set(k, user.UserID)
			continue
		}
		if strings.Contains(k, "password") {
			postData.Set(k, svpn_password)
			continue
		}
	}
	//fmt.Println(postData)

	err = loginVPN(client, postUrl, postData)
	if err != nil {
		return nil, 0, err
	}
	port, statusCode, err := getISP_Port(client)
	if err != nil {
		return nil, statusCode, err
	}
	return fetchISPByVPN(client, port)
}
