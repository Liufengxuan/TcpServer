package com

import (
	"errors"
	"log"

	"github.com/astaxie/beego/config"
)

//IConfig 配置文件接口
type IConfig interface {
	GetInt64(key string) (value int64, err error)
	GetInt(key string) (value int, err error)
	GetString(key string) (value string)
}

var cfg conf

type conf struct {
	configer config.Configer
}

func init() {
	//获取配置文件信息  begin
	log.Println("[正在获取配置文件信息]")
	if cfg.configer == nil {
		c, err := config.NewConfig("ini", "../server.conf")
		if err != nil {
			log.Println("[配置文件初始化失败 ！]")
			return
		}
		if c == nil {
			log.Println("[配置文件无内容 ！]")
			return
		}
		cfg = conf{c}
		log.Println("[已获取到配置文件信息]")
	}
	//获取配置文件信息  end
}

//GetConfig 获取配置文件操作
func GetConfig() (IConfig, error) {
	if cfg.configer == nil {
		return nil, errors.New("获取配置文件失败")
	}
	return &cfg, nil
}

//GetInt64 获取int值
func (config *conf) GetInt64(key string) (value int64, err error) {
	return config.configer.Int64(key)
}

//GetInt 获取int值
func (config *conf) GetInt(key string) (value int, err error) {
	return config.configer.Int(key)
}

//GetString 获取string值
func (config *conf) GetString(key string) (value string) {
	return config.configer.String(key)
}
