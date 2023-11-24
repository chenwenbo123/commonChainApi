package db

import (
	Config "chargeWithdraw/yaml"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb(config *Config.Conf) *gorm.DB {
	dsn := config.Db.UserName + ":" + config.Db.Password + "@tcp(" + config.Db.Host + ":3306)/" + config.Db.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Print(err, "数据库连接错误")
		return nil
	}

	return db
}

func InnerDbInit() *gorm.DB {
	dsn := "yj" + ":" + "8py86dy2Pm6RwhFN" + "@tcp(" + "107.182.187.200" + ":3306)/" + "yj" + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil
	}
	return db
}
