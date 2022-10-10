package ISP

import (
	"fmt"
	"isp_tool/model"
	"net/http"
)

// 开闭原则，添加业务时对扩展开放，对修改关闭，对于扩展是通过接口实现的。
// 如果要增加一个新的打卡人员，只需要实现CheckInTool接口即可。
// 依赖倒置原则，依赖于抽象而不依赖于具体，依赖于接口而不依赖于实现。
type CheckInTool interface {
	CheckIn() error
}

// 单一职责原则，一个类只负责一项职责
type CheckInStudent struct {
	model.Student
	//1表示使用ip地址(若精确到区域则选择IP)，2表示使用isp历史打卡地址(默认)。若配置文件已设置地址，则优先使用配置文件地址
	ChooseLocation int
	Client         *http.Client
}

// 单一职责原则，一个类只负责一项职责
func (cs *CheckInStudent) CheckIn() error {
	var isp_toll ISP_Tool = &CDUStudetForISP{CheckInStudent: *cs}
	err := isp_toll.Login()
	if err != nil {
		return err
	}
	fmt.Println("登录成功")
	err = isp_toll.CheckIn()
	if err != nil {
		return err
	}
	fmt.Println("打卡成功")
	return nil
}
