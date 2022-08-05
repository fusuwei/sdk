package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Databases struct {
	Username string
	Password string
	Host     string
	DBName   string
	Port     int
	DB       *gorm.DB
}

func New(username, password, host, dbName string, port int) *Databases {
	return &Databases{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		DBName:   dbName,
	}
}

func (m *Mysql) OpenConnection() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Addr, m.Port, m.DBName)
	var level gormLogger.LogLevel
	if m.logger.Level == "error" {
		level = gormLogger.Error
	} else {
		level = gormLogger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(m.logger.GetLogger(), gormLogger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  level,
		}),
	})
	if err != nil {
		m.logger.Fatalf(err.Error())
		return
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		m.logger.Fatalf(err.Error())
		return
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	m.DB = db
	go m.reconnection(sqlDB)
	return
}

func (m *Mysql) reconnection(db *sql.DB) {
	for {
		if err := db.Ping(); err != nil {
			time.Sleep(5 * time.Second)
			m.OpenConnection()
			break
		}
		time.Sleep(time.Second * 5)
	}

}

func (m *Mysql) Scopes(option ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	return m.DB.Scopes(option...)
}

func (m *Mysql) DBER() *gorm.DB {
	return m.DB
}

