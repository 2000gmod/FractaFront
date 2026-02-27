package diag

import (
	"fmt"
	"strings"
)

type ErrorContainer struct {
	Message  string
	Filaname string
	Line     int
}

func (e *ErrorContainer) Error() string {
	return fmt.Sprintf("(%s:%d) %s", e.Filaname, e.Line, e.Message)
}

type ErrorList []*ErrorContainer

func (el ErrorList) Error() string {
	sb := strings.Builder{}

	for _, v := range el {
		sb.WriteString(v.Error())
	}
	return sb.String()
}

func CreateError(msg, file string, line int) *ErrorContainer {
	return &ErrorContainer{
		Message:  msg,
		Filaname: file,
		Line:     line,
	}
}

func DiagnoseErrors(list ErrorList) {
	for _, v := range list {
		fmt.Println(v.Error())
	}
}
