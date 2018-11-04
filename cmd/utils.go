package cmd

import (
	"strings"
)

//SplitString 根据','或者';'拆分字符串
func SplitString(str string) (strList []string) {
	if strings.Contains(str, ",") {
		strList = strings.Split(str, ",")
	} else if strings.Contains(str, ";") {
		strList = strings.Split(str, ";")
	} else {
		strList = []string{str}
	}
	return
}
