package service

import (
	"regexp"

	"github.com/fitant/xbin-api/src/types"
)

var regex, _ = regexp.Compile("[A-Z]")

func checkNoteType(s string) types.SnippetType {
	r := regex.FindAllStringIndex(s, -1)
	if r == nil {
		return types.InvalidSnippet
	}
	lr := len(r)
	switch lr {
	case 1:
		return types.StaticSnippet
	case 2:
		return types.EphemeralSnippet
	case 3:
		return types.ProlongedSnippet
	}
	return types.InvalidSnippet
}
