package dao

import (
	"log"
	"os"
	"time"

	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// init 初始化数据库连接并设置默认的 DAO 数据库实例
// 使用当前工作目录下的 db/db.sqlite 作为 SQLite 数据库文件
// 如果连接失败会直接 panic
func init() {
	dbDir := env.DataDir
	common.CreateDir(dbDir)
	dbPath := dbDir + "/db.sqlite"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,          // Don't include params in the SQL log
				Colorful:                  false,         // Disable color
			},
		),
	})
	if err != nil {
		log.Println("数据库初始化失败")
		panic(err)
	}
	DB = db
	log.Println("数据库初始化成功")
	initTask()
}
