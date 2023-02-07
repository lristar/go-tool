package pgs

//
//import (
//	"context"
//	"fmt"
//	"github.com/opentracing/opentracing-go"
//	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
//	"io"
//	configs "platform-backend/config"
//	"strconv"
//
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//)
//
//var db *gorm.DB
//
//type PgConfig struct {
//	PgHost       string `json:"pg_host"`
//	PgPort       string `json:"pg_port"`
//	PgUname      string `json:"pg_uname"`
//	PgPwd        string `json:"pg_pwd"`
//	PgDbName     string `json:"pg_db_name"`
//	MaxOpenConns string `json:"max_open_conns"`
//	MaxIdleConns string `json:"max_idle_conns"`
//}
//
//func (p *PgConfig) Start(serverKey string) (io.Closer, error) {
//	var err error
//	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable TimeZone=Asia/Shanghai", p.PgUname, p.PgPwd, p.PgDbName, p.PgPort, p.PgHost)
//	db, err = gorm.Open(postgres.New(postgres.Config{
//		DSN:                  dsn,
//		PreferSimpleProtocol: true,
//	}), &gorm.Config{})
//	if err != nil {
//		logger.Errorf("pg连接GORM初始化失败1:%s\n", err.Error())
//	}
//	if configs.GetConf("JaegerEnable") == "true" {
//		// 3. 最重要的一步，使用我们定义的插件
//		_ = db.Use(&OpentracingPlugin{})
//		// 4. 生成新的Span - 注意将span结束掉，不然无法发送对应的结果
//		span := opentracing.StartSpan("gormTracing unit test")
//		defer span.Finish()
//
//		// 5. 把生成的Root Span写入到Context上下文，获取一个子Context
//		// 通常在Web项目中，Root Span由中间件生成
//		ctx := opentracing.ContextWithSpan(context.Background(), span)
//
//		// 6. 将上下文传入DB实例，生成Session会话
//		// 这样子就能把这个会话的全部信息反馈给Jaeger
//		db = db.WithContext(ctx)
//	}
//	cons, err := db.DB()
//	if err != nil {
//		logger.Errorf("pg连接GORM初始化失败2:%s\n", err.Error())
//	}
//	var n int
//	if p.MaxOpenConns == "" {
//		n = 10
//	} else {
//		n, _ = strconv.Atoi(p.MaxOpenConns)
//	}
//	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
//	cons.SetMaxIdleConns(p.MaxIdleConns)
//	// SetMaxOpenConns 设置打开数据库连接的最大数量
//	cons.SetMaxOpenConns(n)
//	// SetConnMaxLifetime 设置了连接可复用的最大时间
//	//cons.SetConnMaxLifetime(time.Hour)
//	return cons, nil
//}
//
//// GetConn 获取pg连接 gorm
//func GetConn() *gorm.DB {
//	if configs.GetConf("Debug") == "true" {
//		return db.Debug()
//	}
//	return db
//}
