package mgo

import (
	"time"

	"github.com/imaimang/mgo/logs"

	"github.com/BurntSushi/toml"
)

//反序列化配置文件 文件格式为toml
func UnmarshalConfig(configFile string, config interface{}) error {
	_, err := toml.DecodeFile(configFile, config)
	return err
}

//savePath 文件保存路径 默认 ./logs/
//saveTimeout 文件保存时间 默认为0无限制
//maxFileSize 单个文件存储大小 单位 MB 默认 10MB
func NewLog(savePath string, saveTimeout time.Duration, logFieSize float32) *logs.Logger {
	logger := logs.NewLogger()
	logger.SetSaveOption(savePath, saveTimeout, logFieSize)
	return logger
}
