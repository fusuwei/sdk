package db

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Databases struct {
	Dialect         gorm.Dialector
	DB              *gorm.DB
	sqlDB           *sql.DB
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
}

func New(dialect gorm.Dialector) *Databases {
	return &Databases{
		Dialect: dialect,
	}
}

func (d *Databases) Connection(maxIdleConns, maxOpenConns int, connMaxLifetime time.Duration) error {
	db, err := gorm.Open(d.Dialect)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(maxIdleConns)
	d.maxIdleConns = maxIdleConns
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)
	d.maxOpenConns = maxOpenConns
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	d.connMaxLifetime = connMaxLifetime
	go d.reconnection(sqlDB)
	d.DB = db
	return nil
}

func (d *Databases) reconnection(db *sql.DB) {
	for {
		if err := db.Ping(); err != nil {
			time.Sleep(5 * time.Second)
			d.Connection(d.maxIdleConns, d.maxOpenConns, d.connMaxLifetime)
			break
		}
		time.Sleep(time.Second * 5)
	}

}

func (d *Databases) SetLogger(p gormLogger.Interface) {
	d.DB.Logger = p
}

func (d *Databases) Close() {
	d.sqlDB.Close()
}
