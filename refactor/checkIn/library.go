package checkIn

import "time"

type LibraryCheckInTool interface {
	CheckIn_Library(time.Time) error
}

type LibraryCheckInToolPerson interface {
	UseCheckInTool(tool LibraryCheckInTool) error
}

func CheckIn_Library(p LibraryCheckInToolPerson, tool LibraryCheckInTool) error {
	return p.UseCheckInTool(tool)
}
