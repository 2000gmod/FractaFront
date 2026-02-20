package diag

import (
	"fmt"
	"sync"
)

var GlobalErrors = make([]*ErrorContainer, 0)
var globalErrorsMux = sync.RWMutex{}

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

func AppendError(errs ...*ErrorContainer) {
	globalErrorsMux.Lock()
	defer globalErrorsMux.Unlock()
	GlobalErrors = append(GlobalErrors, errs...)
}

func HadErrors() bool {
	globalErrorsMux.RLock()
	defer globalErrorsMux.RUnlock()
	return len(GlobalErrors) != 0
}

func ReportErrors() {
	globalErrorsMux.RLock()
	defer globalErrorsMux.RUnlock()

	if !HadErrors() {
		return
	}

	for _, e := range GlobalErrors {
		fmt.Printf("%v\n", e)
	}
}
