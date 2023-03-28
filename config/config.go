package config

import (
	"flag"
	"github.com/golang/glog"
	"github.com/jinzhu/configor"
)

type Config struct {
	Addr  Addr  `required:"true" yaml:"Addr"`
	Mysql Mysql `required:"true" yaml:"Mysql"`
}

type Addr struct {
	HttpAddr string `required:"true" yaml:"HttpAddr"` // 提供的http服务地址
	GrpcAddr string `required:"true" yaml:"GrpcAddr"` // 提供的grpc服务地址
}

type Mysql struct {
	DSN             string `required:"true" yaml:"DSN"`           // 数据库连接DSN
	MaxIdleConns    int    `yaml:"MaxIdleConns" default:"5"`      // 最大闲置连接数
	MaxOpenConns    int    `yaml:"MaxOpenConns" default:"10"`     // 最大连接数
	ConnMaxLifetime int    `yaml:"ConnMaxLifetime" default:"180"` // 连接最长生命周期时间
	ConnMaxIdleTime int    `yaml:"ConnMaxIdleTime" default:"60"`  // 连接最长闲置时间
}

var config Config

func init() {
	// 解析flag
	confPath := flag.String("c", "config/config.yaml", "配置文件路径")
	flag.Parse()
	glog.Infoln("开始加载配置文件")
	err := configor.Load(&config, *confPath)
	if err != nil {
		panic(err)
	}
}

func GetConf() *Config {
	return &config
}
