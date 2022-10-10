package model

type Student struct {
	StudentInfo
}

// Info不看做类，不实现接口，只是一个结构体
type StudentInfo struct {
	Name string
	Age  int
	//身份证号
	IdCard     string
	SchoolInfo StudentSchoolInfo
	CurrentLocation
}

type StudentSchoolInfo struct {
	StudentId string //学号
	Class     string
	Year      int    //入学年份
	Major     string //专业
	VpnPwd    string
	ISP_Pwd   string
}

// 当前位置,省市区
type CurrentLocation struct {
	Province string
	City     string
	Area     string
}
