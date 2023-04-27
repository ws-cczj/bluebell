package main

import (
	"bluebell/crontab"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
)

func main() {
	InitDevs()
	// 路由配置
	r := routes.Setup()

	_ = r.Run(settings.Conf.AppConfig.Port)

	mysql.Close()
	redis.Close()
	crontab.NewCrontabInstance().StopAll() // 关闭定时任务
}

func InitDevs() {
	// 初始化数据库连接
	mysql.InitMysql()
	// 初始化redis
	redis.InitRedis()
	// 初始化雪花ID生成器
	snowflake.InitSnowID(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineID)
	// 初始化定时任务
	crontab.NewCrontabInstance()
}
