package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL 进行封装后的 MySQL 客户端，使用前需要先执行"MySQL.InitDefault"或"MySQL.InitUse"函数
type MySQL struct {
	*gorm.DB
	InitError error
}

// InitDefault 根据 dsn 初始化默认配置的 MySQL 客户端
//
// DSN-Template: root:root@tcp(127.0.0.1:3306)/kspace?charset=utf8mb4&parseTime=True&loc=Local
//
// 如果使用 autoMigrate 参数，将对这些模型进行自动迁移
func (slf *MySQL) InitDefault(dsn string, autoMigrate ...any) MySQL {
	slf.DB, slf.InitError = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，GORM 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，GORM 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，GORM 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 GORM 版本自动配置
	}), &gorm.Config{})

	if len(autoMigrate) > 0 {
		slf.InitError = slf.DB.AutoMigrate(autoMigrate...)
	}
	return *slf
}

// InitUse 根据特定配置（mysql.Config 、*gorm.Config）初始化 MySQL 客户端
//
// 如果使用 autoMigrate 参数，将对这些模型进行自动迁移
func (slf *MySQL) InitUse(mysqlConfig mysql.Config, gormConfig *gorm.Config, autoMigrate ...any) MySQL {
	slf.DB, slf.InitError = gorm.Open(mysql.New(mysqlConfig), gormConfig)

	if len(autoMigrate) > 0 {
		slf.InitError = slf.DB.AutoMigrate(autoMigrate...)
	}
	return *slf
}
