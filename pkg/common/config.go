package common

import (
	"go.uber.org/zap/zapcore"
)

// 系统配置, 配置字段可参见yml注释
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
type Configuration struct {
	System SystemConfiguration `mapstructure:"system" json:"system"`
	Logs   LogsConfiguration   `mapstructure:"logs" json:"logs"`
	Mysql  MysqlConfiguration  `mapstructure:"mysql" json:"mysql"`
	Redis  RedisConfiguration  `mapstructure:"redis" json:"redis"`
	Jwt    JwtConfiguration    `mapstructure:"jwt" json:"jwt"`
	Upload UploadConfiguration `mapstructure:"upload" json:"upload"`
	Casbin CasbinConfiguration `mapstructure:"casbin" json:"casbin"`
}

type SystemConfiguration struct {
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	AppMode         string `mapstructure:"app-mode" json:"appMode"`
	Port            int    `mapstructure:"port" json:"port"`
	OperationLogKey string `mapstructure:"operation-log-key" json:"operationLogKey"`
	Key             string `mapstructure:"key" json:"key"`
}

type LogsConfiguration struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type MysqlConfiguration struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}
type RedisConfiguration struct {
	Host     string `mapstructure:"host" json:"host"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
	Port     int    `mapstructure:"port" json:"port"`
}

type CasbinConfiguration struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
	LoadDelay int    `mapstructure:"load-delay" json:"loadDelay"`
}

type JwtConfiguration struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type UploadConfiguration struct {
	SaveDir       string `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize uint   `mapstructure:"single-max-size" json:"singleMaxSize"`
}
