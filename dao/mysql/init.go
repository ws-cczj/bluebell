package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func InitMysql(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		// -- 将日志记录到日志库中
		zap.L().Error("mysql connect is fail", zap.Error(err))
		return err
	}
	db.SetMaxIdleConns(cfg.MaxIdles)
	db.SetMaxOpenConns(cfg.MaxOpens)
	return nil
}

func Close() {
	_ = db.Close()
}
