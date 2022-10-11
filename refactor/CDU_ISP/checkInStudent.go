package CDU_ISP

import (
	"CDU_Tool/checkIn"
	"CDU_Tool/model"
)

// 单一职责原则，一个类只负责一项职责
type CDU_CheckInStudent struct {
	model.CDU_Student
	//1表示使用ip地址(若精确到区域则选择IP)，2表示使用isp历史打卡地址(默认)。若配置文件已设置地址，则优先使用配置文件地址
	ChooseLocation int
}

// 业务逻辑层，人使用工具
func (cs *CDU_CheckInStudent) UseCheckInTool(tool checkIn.CheckInTool) error {
	err := tool.CheckIn()
	if err != nil {
		return err
	}
	return nil
}
