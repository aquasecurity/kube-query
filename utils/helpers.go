package utils

import (
	"fmt"
)

func Bool2str(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func Map2Str(m interface{}) (merged string) {
	for k, v := range m.(map[string]string) {
		merged += fmt.Sprintf("%v=%v,", k, v)
	}
	return
}