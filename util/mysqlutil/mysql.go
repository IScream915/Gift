package mysqlutil

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"` // 单位秒
	TablePrefix  string `mapstructure:"table_prefix"`
	LogLevel     string `mapstructure:"log_level"`
}

func NewMysqlClient(ctx context.Context) (*gorm.DB, error) {
	config := &Config{}

	// 载入dbconfig
	if err := loadDbConfig(config); err != nil {
		return nil, err
	}

	// 填充默认值
	repairClient(config)

	// 启用数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)

	client, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,  // string 类型字段的默认长度
			DisableDatetimePrecision:  true, // 禁用 datetime 精度， MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true, // 重命名索引时采用删除并创建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false,
		}),
		&gorm.Config{
			PrepareStmt: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: config.TablePrefix,
			},
			//Logger: log.NewSqlLogger(getLogLevel(config.LogLevel), false, 200*time.Millisecond),
		},
	)

	sqlDB, err := client.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	return client, nil
}

// 修复数据库连接
func repairClient(c *Config) {
	if c.Host == "" {
		c.Host = DefaultHost
	}
	if c.Port <= 0 {
		c.Port = DefaultPort
	}
	if c.User == "" {
		c.User = DefaultUser
	}
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = DefaultMaxOpenConns
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = DefaultMaxIdleConns
	}
	if c.MaxLifetime <= 0 {
		c.MaxLifetime = DefaultMaxLifetime
	}
	if c.TablePrefix == "" {
		c.TablePrefix = DefaultTablePrefix
	}
	if c.LogLevel == "" {
		c.LogLevel = DefaultLogLevel
	}
}

func loadDbConfig(c *Config) error {
	// 读取文件内容
	configPath := DefaultDbConfigPath

	// 使用 Viper 加载配置文件
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	// 获取配置文件内容
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	return nil
}
