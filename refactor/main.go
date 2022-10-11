package main

import (
	"CDU_Tool/CDU_ISP"
	"CDU_Tool/checkIn"
	"fmt"
	"net/http"
)

func main() {
	stu := CDU_ISP.CDU_CheckInStudent{}
	stu.Name = "张三"
	stu.University = "Chengdu University"
	stu.StudentId = "201800000000"
	stu.VpnPwd = "123456"
	stu.Age = 18
	var stu1 checkIn.CheckInToolPerson = &stu
	isptool := CDU_ISP.ISP_Tool{
		Stu:    &stu,
		Client: &http.Client{},
	}
	err := stu1.UseCheckInTool(&isptool)
	if err != nil {
		fmt.Println(err)
		return
	}
}
