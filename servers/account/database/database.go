package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func CreateMysqlClient(dsn string) (*gorm.DB, error) {
	dbc, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableNestedTransaction: true, //关闭嵌套事务, Nested : adj, 嵌套的
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Logger: log.NewGormLogrus(),
	})
	logger.Default = logger.Default.LogMode(logger.Info)
	return dbc, err
}
