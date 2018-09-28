package com

import (
	"crypto/md5"

	"encoding/hex"
)

//GetMD5 根据输入获取加密后的md5
func GetMD5(str string) (string, error) {

	md5Ctx := md5.New()

	_, err := md5Ctx.Write([]byte(str))
	if err != nil {
		return "", err
	}

	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr), nil

}
