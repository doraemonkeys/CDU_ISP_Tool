package checkIn

// 开闭原则，添加业务时对扩展开放，对修改关闭，对于扩展是通过接口实现的。
// 如果要增加一个新的打卡人员，只需要实现CheckInTool接口即可。
// 依赖倒置原则，依赖于抽象而不依赖于具体，依赖于接口而不依赖于实现。
// 抽象层
type CheckInTool interface {
	CheckIn() error
}

// 业务逻辑层抽象
type CheckInToolPerson interface {
	UseCheckInTool(tool CheckInTool) error
}
