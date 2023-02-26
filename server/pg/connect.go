package pgs

import (
	"fmt"
	"github.com/lristar/go-tool/server/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

const (
	DEFAULTMAXOPENCONN = 10
	DEFAULTMAXIDlECONN = 10
)

var (
	db   *PgClient
	once sync.Once
)

type Config struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Uname       string `json:"uname"`
	Pwd         string `json:"pwd"`
	Dbname      string `json:"dbname"`
	MaxIdleConn int    `json:"max_idle_conn"`
	MaxOpenConn int    `json:"max_open_conn"`
}

type PgClient struct {
	*gorm.DB
}

func InitPg(c Config) {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s sslmode=disable TimeZone=Asia/Shanghai", c.Uname, c.Pwd, c.Dbname, c.Port, c.Host)
		conn, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			logger.Errorf("pg连接GORM初始化失败1:%s\n", err.Error())
		}

		cons, err := conn.DB()
		if err != nil {
			logger.Errorf("pg连接GORM初始化失败2:%s\n", err.Error())
		}
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		mic := c.MaxIdleConn
		if c.MaxOpenConn == 0 {
			mic = DEFAULTMAXIDlECONN
		}
		cons.SetMaxIdleConns(mic)
		// SetMaxOpenConns 设置打开数据库连接的最大数量
		n := c.MaxOpenConn
		if c.MaxOpenConn == 0 {
			n = DEFAULTMAXOPENCONN
		}
		cons.SetMaxOpenConns(n)
		// SetConnMaxLifetime 设置了连接可复用的最大时间
		cons.SetConnMaxLifetime(time.Hour)
		db = &PgClient{conn}
		logger.Infof("pg:host->%s,db_name->%s ---\n", c.Host, c.Dbname)
	})
}

// GetDb 获取pg连接 gorm
func GetDb(isDebug bool) *PgClient {
	if isDebug {
		return &PgClient{db.Debug()}
	}
	return db
}

func (p *PgClient) Close() error {
	return nil
}
