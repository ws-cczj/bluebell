package logger

import (
	"bluebell/settings"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化日志
func InitLogger(cfg *settings.LogConfig) error {
	encoder := getEncode()
	syncer := getLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxAge,
		cfg.MaxBackups,
	)
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	// -- 用于开发者模式将日志打印到控制台
	var core zapcore.Core
	if settings.Conf.AppConfig.Mode == gin.DebugMode {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, syncer, l),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l),
		)
	} else {
		core = zapcore.NewCore(encoder, syncer, l)
	}

	lg := zap.New(core, zap.AddCaller()) // --添加函数调用信息
	zap.ReplaceGlobals(lg)               // -- 将日志替换到全局的zaplogger中去
	return nil
}

// getEncode 进行自定义的日志输出格式
func getEncode() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// -- 自定义时间格式
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     getEncodeTime(),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter 进行日志切割
func getLogWriter(filename string, max_size, max_age, max_backups int) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,    // -- 日志文件名
		MaxSize:    max_size,    // -- 最大日志数 M为单位!!!
		MaxAge:     max_age,     // -- 最大存在天数
		MaxBackups: max_backups, // -- 最大备份数量
		Compress:   false,       // --是否压缩
	}
	return zapcore.AddSync(lumberjackLogger)
}

// encodeTimeLayout 进行时间格式的解析
func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

// getEncodeTime 自定义的日志打印时间格式
func getEncodeTime() func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		encodeTimeLayout(t, "[2006-01-02 15:04:05]", enc)
	}
}
