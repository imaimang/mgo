package dblib

import (
	"errors"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//NewMysql NewMysql
//user:pwd@tcp(ip:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func NewMysql(dbAddress string) (*gorm.DB, error) {

	indexEnd := strings.Index(dbAddress, "?")
	if indexEnd == -1 {
		indexEnd = len(dbAddress)
	}
	indexStart := strings.Index(dbAddress, "/")
	if indexStart == -1 {
		return nil, errors.New("db address foramt faild")
	}
	dbName := dbAddress[indexStart+1 : indexEnd]
	dbAddressNoDB := dbAddress[0:indexStart+1] + dbAddress[indexEnd:]

	db, err := gorm.Open(mysql.Open(dbAddressNoDB), &gorm.Config{})
	if err == nil {
		err = db.Exec("CREATE DATABASE  IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8 COLLATE utf8_bin").Error
		if err != nil {
			return nil, err
		}
		dbConnect, err := db.DB()
		if err != nil {
			return nil, err
		}
		dbConnect.Close()
		db, err = gorm.Open(mysql.Open(dbAddress), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			SkipDefaultTransaction: true})
		if err == nil {
			dbConnect, err = db.DB()
			if err != nil {
				return nil, err
			}
			dbConnect.SetMaxIdleConns(10)
			dbConnect.SetMaxOpenConns(100)
		}
	}
	return db, err
}
