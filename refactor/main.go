package main

import (
	"CDU_Tool/CDU_ISP/student"
	"CDU_Tool/CDU_ISP/teacher"
	"CDU_Tool/checkIn"
	"fmt"
	"net/http"
)

func main() {
	stu := student.CDU_CheckInStudent{}
	stu.Name = "张三"
	stu.University = "Chengdu University"
	stu.StudentId = "201800000000"
	stu.VpnPwd = "123456"
	stu.Age = 18

	isptool := student.Stu_ISP_Tool{
		Client: &http.Client{},
	}

	err := checkIn.CheckIn(&stu, &isptool)
	if err != nil {
		fmt.Println(err)
		return
	}

	teacher1 := teacher.CDU_CheckInTeacher{}
	teacher1.Name = "李四"
	teacher1.Age = 30
	teacher1.TeacherId = "201809000000"
	teacher1.VpnPwd = "123456"
	teacher1.ISP_Pwd = "123456"

	isptool2 := teacher.Teacher_ISP_Tool{
		Client: &http.Client{},
	}
	err = checkIn.CheckIn(&teacher1, &isptool2)
	if err != nil {
		fmt.Println(err)
		return
	}
}
