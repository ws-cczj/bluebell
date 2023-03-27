package settings

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Config) // -- 定义全局变量进行存储配置信息,采用new是为了确保指针地址不发生变化

type Config struct {
	AppConfig `mapstructure:"app"`
	LogConfig `mapstructure:"log"`
	Mdb       MysqlConfig `mapstructure:"mysql"`
	Rdb       RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Mode      string `mapstructure:"mode"`
	Port      string `mapstructure:"port"`
	Jwt       `mapstructure:"jwt"`
	SnowFlake `mapstructure:"snowflake"`
	RateLimit `mapstructure:"ratelimit"`
}

type Jwt struct {
	AtokenAt int64 `mapstructure:"atoken_at"`
	RtokenAt int64 `mapstructure:"rtoken_at"`
}

type SnowFlake struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type RateLimit struct {
	GenInterval int64 `mapstructure:"gen_interval"`
	MaxCaps     int64 `mapstructure:"max_caps"`
}

type LogConfig struct {
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Level      string `mapstructure:"level"`
	Layout     string `mapstructure:"layout"`
	Filename   string `mapstructure:"filename"`
}

type MysqlConfig struct {
	MaxIdles int    `mapstructure:"max_idles_conns"`
	MaxOpens int    `mapstructure:"max_opens_conns"`
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Dbname   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

func init() {
	viper.SetConfigFile("./conf/config.json")
	//viper.SetConfigName("config")  // --设置配置文件得名称(不能获取后缀)
	//viper.SetConfigType("yaml")    // -- 设置配置文件的类型(专用于远程获取配置信息时进行使用)
	//viper.AddConfigPath("./conf/") // -- 设置配置文件的路径
	// 统一处理错误
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}()
	err = viper.ReadInConfig()

	// -- 将配置信息反序列化到 Conf 全局变量中去
	err = viper.Unmarshal(Conf)

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Debug("config was changed!", zap.String("name", in.Name))
		if err = viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	})
}
