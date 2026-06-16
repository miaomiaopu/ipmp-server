package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Encryption EncryptionConfig
	CORS       CORSConfig
	AI         AIConfig
	Log        LogConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Type     string         `mapstructure:"type"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
}

func (d DatabaseConfig) DSN() string {
	if d.Type == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			d.MySQL.User, d.MySQL.Password, d.MySQL.Host, d.MySQL.Port, d.MySQL.Name, d.MySQL.Charset)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Postgres.Host, d.Postgres.Port, d.Postgres.User, d.Postgres.Password, d.Postgres.Name, d.Postgres.SSLMode)
}

func (d DatabaseConfig) Driver() string {
	if d.Type == "mysql" {
		return "mysql"
	}
	return "postgres"
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type CORSConfig struct {
	Origins string `mapstructure:"origins"`
}

func (c CORSConfig) AllowedOrigins() []string {
	return strings.Split(c.Origins, ",")
}

type AIConfig struct {
	Provider string           `mapstructure:"provider"`
	DeepSeek AIProviderConfig `mapstructure:"deepseek"`
	OpenAI   AIProviderConfig `mapstructure:"openai"`
	Claude   AIProviderConfig `mapstructure:"claude"`
}

type AIProviderConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 绑定环境变量默认值 (BindEnv 仅对空字符串出错，常量 key 无此风险)
	_ = v.BindEnv("server.port", "SERVER_PORT")
	_ = v.BindEnv("server.mode", "SERVER_MODE")

	_ = v.BindEnv("database.type", "DB_TYPE")
	_ = v.BindEnv("database.postgres.host", "PG_HOST")
	_ = v.BindEnv("database.postgres.port", "PG_PORT")
	_ = v.BindEnv("database.postgres.user", "PG_USER")
	_ = v.BindEnv("database.postgres.password", "PG_PASSWORD")
	_ = v.BindEnv("database.postgres.name", "PG_NAME")
	_ = v.BindEnv("database.postgres.sslmode", "PG_SSLMODE")
	_ = v.BindEnv("database.mysql.host", "MYSQL_HOST")
	_ = v.BindEnv("database.mysql.port", "MYSQL_PORT")
	_ = v.BindEnv("database.mysql.user", "MYSQL_USER")
	_ = v.BindEnv("database.mysql.password", "MYSQL_PASSWORD")
	_ = v.BindEnv("database.mysql.name", "MYSQL_NAME")
	_ = v.BindEnv("database.mysql.charset", "MYSQL_CHARSET")

	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	_ = v.BindEnv("encryption.key", "IPMP_ENCRYPTION_KEY")
	_ = v.BindEnv("cors.origins", "CORS_ORIGINS")
	_ = v.BindEnv("ai.deepseek.api_key", "DEEPSEEK_API_KEY")
	_ = v.BindEnv("ai.openai.api_key", "OPENAI_API_KEY")
	_ = v.BindEnv("ai.claude.api_key", "CLAUDE_API_KEY")
	_ = v.BindEnv("log.level", "LOG_LEVEL")
	_ = v.BindEnv("log.format", "LOG_FORMAT")

	// 默认值
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.postgres.host", "localhost")
	v.SetDefault("database.postgres.port", 5432)
	v.SetDefault("database.postgres.user", "postgres")
	v.SetDefault("database.postgres.name", "ipmp")
	v.SetDefault("database.postgres.sslmode", "disable")
	v.SetDefault("database.mysql.host", "localhost")
	v.SetDefault("database.mysql.port", 3306)
	v.SetDefault("database.mysql.user", "root")
	v.SetDefault("database.mysql.name", "ipmp")
	v.SetDefault("database.mysql.charset", "utf8mb4")
	v.SetDefault("jwt.access_expire", "24h")
	v.SetDefault("jwt.refresh_expire", "168h")
	v.SetDefault("cors.origins", "http://localhost:5173")
	v.SetDefault("ai.provider", "mock")
	v.SetDefault("ai.deepseek.model", "deepseek-v4-flash")
	v.SetDefault("ai.deepseek.base_url", "https://api.deepseek.com")
	v.SetDefault("ai.openai.model", "gpt-4o")
	v.SetDefault("ai.claude.model", "claude-sonnet-4-5")
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "console")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 解析 duration 字符串
	accessExpire, err := time.ParseDuration(v.GetString("jwt.access_expire"))
	if err != nil {
		return nil, fmt.Errorf("parse jwt.access_expire: %w", err)
	}
	cfg.JWT.AccessExpire = accessExpire

	refreshExpire, err := time.ParseDuration(v.GetString("jwt.refresh_expire"))
	if err != nil {
		return nil, fmt.Errorf("parse jwt.refresh_expire: %w", err)
	}
	cfg.JWT.RefreshExpire = refreshExpire

	return &cfg, nil
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read .env: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set .env key %s: %w", key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan .env: %w", err)
	}
	return nil
}
