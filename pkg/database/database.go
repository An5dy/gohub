package database

import (
	"database/sql"
	"errors"
	"fmt"
	"gohub/pkg/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 对象
var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector, _logger logger.Interface) {
	// 使用 gorm.Open 连接数据库
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})

	// 错误处理
	if err != nil {
		fmt.Println(err.Error())
	}

	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 获取当前数据库名
func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() error {
	var err error

	switch config.Get("database.connection") {
	case "mysql":
		deleteMysqlDatabase()
	case "sqlite":
		deleteAllSqliteTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}

// 删除所有 sqlite 中的表
func deleteAllSqliteTables() error {
	tables := []string{}
	DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'")
	for _, table := range tables {
		DB.Migrator().DropTable(table)
	}
	return nil
}

// 删除所有 mysql 中当前数据库的表
func deleteMysqlDatabase() error {
	dbname := CurrentDatabase()
	sql := fmt.Sprintf("DROP DATABASE %s;", dbname)
	if err := DB.Exec(sql).Error; err != nil {
		return err
	}

	sql = fmt.Sprintf("CREATE DATABASE %s;", dbname)
	if err := DB.Exec(sql).Error; err != nil {
		return err
	}
	sql = fmt.Sprintf("USE %s;", dbname)
	if err := DB.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
