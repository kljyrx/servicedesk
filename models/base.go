package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/os/gfile"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

// 在其它model的实体类中可直接调用
var db *gorm.DB

func InitDB() *gorm.DB {
	var err error

	// sqlite3配置信息
	sqliteName := "serviceDesk"

	dataSource := "database" + string(os.PathSeparator) + sqliteName
	if !gfile.Exists(dataSource) {
		os.MkdirAll(path.Dir(dataSource), os.ModePerm)
		os.Create(dataSource)
	}
	db, err = gorm.Open("sqlite3", dataSource)

	if err != nil {
		db.Close()
	}

	// 设置连接池，空闲连接
	db.DB().SetMaxIdleConns(50)
	// 打开链接
	db.DB().SetMaxOpenConns(100)

	// 表明禁用后缀加s
	db.SingularTable(true)

	db.LogMode(true)

	return db
}
