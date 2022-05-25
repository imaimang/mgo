package mgo

import (
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/imaimang/mgo/dblib"
	"github.com/imaimang/mgo/httplib"
	"github.com/imaimang/mgo/logs"
	"gorm.io/gorm"
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

func NewHttpServer() *httplib.HTTPServer {
	httpServer := new(httplib.HTTPServer)
	return httpServer
}

func NewHttpProxy(ip string, port int) *httplib.HttpProxy {
	httpProxy := new(httplib.HttpProxy)
	httpProxy.SetBaseUrl("http://" + ip + ":" + strconv.Itoa(port))
	return httpProxy
}

//NewMysql NewMysql
//user:pwd@tcp(ip:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func NewMysql(dbAddress string) (*gorm.DB, error) {
	return dblib.NewMysql(dbAddress)
}
