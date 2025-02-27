package nodeinfodb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(dbPath string) {
	// 连接到 SQLite 数据库
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// 自动迁移模型，确保表存在
	db.AutoMigrate(&NodeInfo{})
}

func GetDB() *gorm.DB {
	return db
}
