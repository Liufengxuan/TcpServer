package com

import (
	"strings"
)

func CmdFormat(s string) []string {
	var result []string
	//去除回车和换行
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)

	//按空格分割为数组。
	sArr := strings.Split(s, " ")
	length := len(sArr)
	//遍历每个元素、去除多余空格。
	for i := 0; i < length; i++ {
		if sArr[i] == " " || sArr[i] == "" {
			continue
		}
		sArr[i] = strings.Replace(sArr[i], " ", "", -1)

		if sArr[i] != "" {
			result = append(result, sArr[i])
		}
	}

	return result
}
