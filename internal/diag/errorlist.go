package diag

import "strings"

type ErrorList struct {
	Errors []*ErrorContainer
}

func (es *ErrorList) Error() string {
	bs := strings.Builder{}

	bs.WriteString("The following errors ocurred:\n")

	for _, v := range es.Errors {
		bs.WriteString(v.Error())
		bs.WriteRune('\n')
	}

	return bs.String()
}
