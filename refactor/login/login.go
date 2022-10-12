package login

type LoginTool interface {
	Login(LoginToolPerson) error
}

type LoginToolPerson interface {
	UseLoginTool(LoginTool) error
}
