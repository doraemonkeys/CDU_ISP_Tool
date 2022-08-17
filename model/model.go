package model

//正则表达式

const VerificationCodeRe string = `"验证码"[^0-9]+([0-9]+)`

const User_no_RE string = `<[^?<>]+\?user_no=([0-9a-z]+)[^>]+>`

const Isp_history_location_Re string = `>([^|]{0,10})\|([^|]+)\|([^<]+)<`

const Ip_locationRe string = `"ct":\"([^"]+)", "prov":"([^"]+)", "city":"([^"]+)", "area":"([^"]+)", "idc":"", "yunyin":"([^"]+)"[^}]+}`

const UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"

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
