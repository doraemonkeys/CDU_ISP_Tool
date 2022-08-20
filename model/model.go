package model

//正则表达式

//"验证码">&nbsp;&nbsp;4707
const VerificationCodeRe string = `"验证码"[^0-9]+([0-9]+)`

//<a onclick="changeIframe('stushow.asp?user_no=996a718272240eeccca0c816704070224652ub1adbd7')">
const User_no_RE string = `<[^?<>]+\?user_no=([0-9a-z]+)[^>]+>`

//>四川省|成都市|龙泉驿区<
const Isp_history_location_Re string = `>([^|]{0,10})\|([^|]+)\|([^<]+)<`

//"ct":"中国", "prov":"四川省", "city":"成都市", "area":"龙泉驿区", "idc":"", "yunyin":"电信", "net":""}
const Ip_locationRe string = `"ct":\"([^"]+)", "prov":"([^"]+)", "city":"([^"]+)", "area":"([^"]+)", "idc":"", "yunyin":"[^"]+"[^}]+}`

const UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"

//程序已经设置自启动
var Auto_Start bool = false

//程序启动前今天的打卡已经完成
var Auto_Clock_IN_Success bool = false

//用户是否使用本程序修改了用户账号配置
var UserConfigChanged bool = false

type UserInfo struct {
	UserName  string
	UserID    string
	UserPwd   string
	UserNonce string //对应isp的user_no字段
	Location
}

type Location struct {
	Province string
	City     string
	Area     string
}
