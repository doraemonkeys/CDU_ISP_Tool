package model

//正则表达式

//"验证码">&nbsp;&nbsp;4707
//const VerificationCodeRe string = `"验证码"[^0-9]+([0-9]+)`

//<a onclick="changeIframe('stushow.asp?user_no=996a718272240eeccca0c816704070224652ub1adbd7')">
//const User_no_RE string = `<[^?<>]+\?user_no=([0-9a-z]+)[^>]+>`

//>四川省|成都市|龙泉驿区<
//const Isp_history_location_Re string = `>([^|]{0,10})\|([^|]+)\|([^<]+)<`

//"ct":"中国", "prov":"四川省", "city":"成都市", "area":"龙泉驿区", "idc":"", "yunyin":"电信", "net":""}
// const Ip_locationRe string = `"ct":\"([^"]+)", "prov":"([^"]+)", "city":"([^"]+)", "area":"([^"]+)", "idc":"", "yunyin":"[^"]+"[^}]+}`

// const Ipv6Re string = `\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*`

const UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"

// const Clock_IN_ID string = `projecthealth_del.asp\?([a-z]+[^=])=([0-9a-z]+)"`

// const Today_statusRe string = `主动登记[^主]+主动登记`

//程序已经设置自启动
var Auto_Start bool = false

//程序启动前今天的打卡已经完成
var Auto_Clock_IN_Success bool = false

//用户是否使用本程序修改了用户账号配置
var UserConfigChanged bool = false

//版本号
var Version string = "v1.5.3"

type Update struct {
	LanzouUrl                 string //蓝奏云地址
	MainProgramDirectUrl      string //直接下载地址
	AutoStartProgramDirectUrl string //直接下载地址
	LanzouPwd                 string //下载密码
	MainProgramName           string //主程序名
	AutoStartProgramName      string //自启动程序名
	AutoStartProgramVersion   string //自启动程序版本
	MainProgramVersion        string //主程序版本
	AutoStartProgramMd5       string //自启动程序md5
	MainProgramMd5            string //主程序md5
}

//全局配置
var All = AllMsg{}

type Header struct {
	Method       string
	Authority    string
	Content_type string
	User_agent   string
	Referer      string
}

type LoginMsg struct {
	LoginWebUrl    string
	LoginUrl       string
	Input1Field    string //学号对应字段
	Input2Field    string //密码对应字段
	Input3Field    string //验证码对应字段
	MatchParamesRe string //匹配参数的正则表达式
	Other          []FieldAndValue
	Head           Header
}

type UserInfo struct {
	UserName  string
	UserID    string
	UserPwd   string
	UserNonce string //对应isp的user_no字段
	Location
	//1表示使用ip地址(默认)，2表示使用isp历史打卡地址。若配置文件已设置地址，则优先使用配置文件地址
	ChooseLocation int
}

type Location struct {
	Province string
	City     string
	Area     string
}

//key-value
type FieldAndValue struct {
	Field string
	Value string
}

type ISPHomeMsg struct {
	ISPHomeUrl string
	Head       Header
}

type ClockInHomeMsg struct {
	ClockInHomeUrl string
	QueryField     string
	Head           Header
}

type ClockInMsg struct {
	ClockInUrl string
	//QueryField    []string
	Head           Header
	ProvinceField  string
	CityField      string
	AreaField      string
	MatchActionRe  string
	MatchParamesRe string
	Other          []FieldAndValue
}

type CancelMsg struct {
	CancelUrl string
	Head      Header
}

type RegexpStr struct {
	VerificationCodeRe      string //验证码
	User_no_Re              string //用户识别码
	Ip_locationRe           string
	Isp_history_location_Re string
	Ipv6Re                  string
	Clock_IN_ID_Re          string
	Today_statusRe          string
	PwdErrorRe              string
	IsNotStudentRe          string //非在校学生
	AbnormalRe              string //打卡异常检测
	AbnormalColorRe         string //打卡异常颜色检测
	Clock_In_success_Re     string //打卡成功
	Already_Clock_In_Re     string //已打卡
	CancelSuccessRe         string //取消打卡成功
	ASN_Home                string //ASN归属地
}

type AllMsg struct {
	User        UserInfo
	Login       LoginMsg
	ISPHome     ISPHomeMsg
	ClockInHome ClockInHomeMsg
	ClockIn     ClockInMsg
	Cancel      CancelMsg
	Regexp      RegexpStr //正则表达式
}
