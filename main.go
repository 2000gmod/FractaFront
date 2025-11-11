//go:generate go run ./gen/gen_errors.go

package main

import (
	"fracta/internal/cmd"
)

func main() {
	cmd.ProgramEntry()
}
