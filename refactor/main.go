package main

import (
	"CDU_Tool/CDU_ISP"
	"fmt"
	"net/http"
)

func main() {
	var stu1 = CDU_ISP.CDU_CheckInStudent{}
	stu1.Name = "张三"
	stu1.Age = 18
	stu1.SchoolInfo.StudentId = "201800000000"
	isptool := CDU_ISP.ISP_Tool{
		Stu:    &stu1,
		Client: &http.Client{},
	}
	err := stu1.UseCheckInTool(&isptool)
	if err != nil {
		fmt.Println(err)
		return
	}
}
