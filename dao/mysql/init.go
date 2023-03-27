package mysql

import (
	"bluebell/settings"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db                 *sqlx.DB
	ErrNoRows          = sql.ErrNoRows
	ErrorUserExist     = errors.New("用户名已经存在")
	ErrorNotComparePwd = errors.New("用户密码不匹配")
)

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		settings.Conf.Mdb.Username,
		settings.Conf.Mdb.Password,
		settings.Conf.Mdb.Host,
		settings.Conf.Mdb.Port,
		settings.Conf.Mdb.Dbname)
	var err error
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		panic(fmt.Errorf("mysql connect fail, err: %s", err))
	}
	db.SetMaxIdleConns(settings.Conf.Mdb.MaxIdles)
	db.SetMaxOpenConns(settings.Conf.Mdb.MaxOpens)
}

func Close() {
	_ = db.Close()
}
