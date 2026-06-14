package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT   JWTConfig   `mapstructure:"jwt"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	ExpireMinutes int    `mapstructure:"expire_minutes"`
}

func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// 环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// 叠加环境专属配置
	env := v.GetString("app.env")
	if env == "" {
		env = "dev"
	}
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.MergeInConfig(); err != nil {
		// 非致命错误——基础配置可能已足够
		fmt.Printf("warn: could not load config.%s.yaml: %v\n", env, err)
	}

	// 环境变量覆盖
	if v := viper.GetString("MYSQL_HOST"); v != "" {
		viper.Set("mysql.host", v)
	}
	if v := viper.GetString("MYSQL_PORT"); v != "" {
		viper.Set("mysql.port", v)
	}
	if v := viper.GetString("MYSQL_DATABASE"); v != "" {
		viper.Set("mysql.database", v)
	}
	if v := viper.GetString("MYSQL_USER"); v != "" {
		viper.Set("mysql.user", v)
	}
	if v := viper.GetString("MYSQL_PASSWORD"); v != "" {
		viper.Set("mysql.password", v)
	}
	if v := viper.GetString("REDIS_HOST"); v != "" {
		viper.Set("redis.host", v)
	}
	if v := viper.GetString("REDIS_PORT"); v != "" {
		viper.Set("redis.port", v)
	}
	if v := viper.GetString("JWT_SECRET"); v != "" {
		viper.Set("jwt.secret", v)
	}
	if v := viper.GetString("JWT_EXPIRE_MINUTES"); v != "" {
		viper.Set("jwt.expire_minutes", v)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 使用主 viper 实例应用环境变量覆盖
	if h := viper.GetString("MYSQL_HOST"); h != "" {
		cfg.MySQL.Host = h
	}
	if p := viper.GetString("REDIS_HOST"); p != "" {
		cfg.Redis.Host = p
	}

	return cfg, nil
}
