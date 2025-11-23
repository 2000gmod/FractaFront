//go:generate go run ./gen/errors/gen_errors.go

package main

import (
	"fracta/cmd"
)

func main() {
	cmd.ProgramEntry()
}
