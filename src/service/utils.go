package service

import (
	"fmt"
	"regexp"
)

func checkIfEphemeral(s string) bool {
	x, _ := regexp.Compile("[A-Z]")
	r := x.FindAllStringIndex(s, -1)
	if r == nil {
		fmt.Errorf("invalid ID")
	}
	return len(r) <= 2
}
