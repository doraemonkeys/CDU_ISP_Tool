package main

import (
	"isp_tool/ISP"
	"isp_tool/model"
)

func main() {
	stu := model.Student{}
	checkInStu := &ISP.CheckInStudent{Student: stu}
	err := checkInStu.CheckIn()
	if err != nil {
		panic(err)
	}
}
