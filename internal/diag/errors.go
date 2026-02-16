package diag

import "fmt"

var GlobalErrors = make([]*ErrorContainer, 0)

type ErrorContainer struct {
	Message  string
	Filaname string
	Line     int
}

func (e *ErrorContainer) Error() string {
	return fmt.Sprintf("(%s:%d) %s", e.Filaname, e.Line, e.Message)
}

func CreateError(msg, file string, line int) *ErrorContainer {
	return &ErrorContainer{
		Message:  msg,
		Filaname: file,
		Line:     line,
	}
}

func AppendError(err *ErrorContainer) {
	GlobalErrors = append(GlobalErrors, err)
}

func HadErrors() bool {
	return len(GlobalErrors) != 0
}

func ReportErrors() {
	if !HadErrors() {
		return
	}
	//fmt.Printf("Got compilation errors\n")
	for _, e := range GlobalErrors {
		fmt.Printf("%v\n", e)
	}
}
