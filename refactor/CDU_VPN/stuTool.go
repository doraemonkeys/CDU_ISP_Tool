package CDU_VPN

import (
	"CDU_Tool/checkIn"
	"CDU_Tool/model"
	"fmt"
	"time"
)

type Stu_VPN_Tool struct {
	BasicTool *Basic_VPN_Tool
	Stu       *VPN_Student
}

type VPN_Student struct {
	model.CDU_Student
	Login *VPN_Login_Person
}

type LibraryCheckInStudent struct {
	*VPN_Student
}

// 图书馆签到
func (tool *Stu_VPN_Tool) CheckIn_Library(checkInTime time.Time) error {
	fmt.Println("图书馆签到时间：", checkInTime.String())
	return nil
}

func (stu *LibraryCheckInStudent) UseCheckInTool(tool checkIn.LibraryCheckInTool) error {
	err := tool.CheckIn_Library(time.Now())
	if err != nil {
		return err
	}
	return nil
}
