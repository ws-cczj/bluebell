package main

import (
	"bluebell/crontab"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化配置文件
	err := settings.InitConfig()
	if err != nil {
		zap.L().Error("init settings fail!", zap.Error(err))
		return
	}
	// 2. 初始化日志
	err = logger.InitLogger(settings.Conf.LogConfig)
	if err != nil {
		zap.L().Error("init logger fail!", zap.Error(err))
		return
	}
	// 3. 初始化数据库连接
	err = mysql.InitMysql(settings.Conf.MysqlConfig)
	if err != nil {
		zap.L().Error("init mysql fail!", zap.Error(err))
		return
	}
	defer mysql.Close()
	err = redis.InitRedis(settings.Conf.RedisConfig)
	if err != nil {
		zap.L().Error("init redis fail!", zap.Error(err))
		return
	}
	defer redis.Close()
	// 4. 初始化分布式ID生成器
	err = snowflake.InitSnowID(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineID)
	if err != nil {
		zap.L().Error("init snowflake fail!", zap.Error(err))
		return
	}
	// 6. 初始化路由
	r := routes.Setup(settings.Conf.AppConfig)
	// 5. 开启定时任务
	ctab := crontab.NewCrontabInstance()
	crontab.MysqlTask()
	crontab.MonitorTask()
	ctab.RunAll()
	// 7. 开启web监听服务,设置优雅关机
	srv := &http.Server{
		Addr:    settings.Conf.AppConfig.Port,
		Handler: r,
	}

	// -- 开启协程进行监听
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen is fail!", zap.Error(err))
		}
	}()

	// -- 等待有中断信号来触发管道信号
	quit := make(chan os.Signal, 1) // -- 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号, 也就是常用的 ctrl + C 终止命令
	// kill -9 发送 syscall.SIGKILL 信号, 但是不能被捕获到， 所以不需要添加该信号
	// signal.Notify会把收到的 syscall.SIGINT 或者 syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit         // -- 如果接收不到信号就在这里一直堵塞
	ctab.StopAll() // 关闭定时任务
	zap.L().Info("Shutdown Server...")
	// -- 创建一个超过5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// -- 5s内优雅关闭服务（将正在处理的服务处理完毕后结束进程）, 超过5s就强制结束
	if err = srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown fail!", zap.Error(err))
	}
	zap.L().Info("Server exiting...")
}
