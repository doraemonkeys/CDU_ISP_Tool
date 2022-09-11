package model

//版本号,与主版本号不相同
var Version string = "1.0.3"

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

const UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"
