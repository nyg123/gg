package dao

import (
	"context"
	"gg/config"
	"github.com/golang/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var once sync.Once
var openDb *gorm.DB

type DB struct {
}

func NewDb() *DB {
	return &DB{}
}

// Conn 获取db连接
func (d *DB) Conn(ctx context.Context) *gorm.DB {
	glog.Infof("db conn: %v", openDb)
	return openDb.WithContext(ctx)
}

func Init(conf config.Mysql) {
	once.Do(func() {
		// init db
		var err error
		openDb, err = gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // 定义logger
		})
		if err != nil {
			glog.Errorf("init db error: %v", err)
			panic(err)
		}
		db, err := openDb.DB()
		if err != nil {
			glog.Errorf("init db error: %v", err)
			panic(err)
		}
		// 设置连接池
		db.SetMaxIdleConns(conf.MaxIdleConns)
		db.SetMaxOpenConns(conf.MaxOpenConns)
		db.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTime) * time.Second)
		glog.Infoln("init db success")
	})
}
