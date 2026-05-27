package main

import (
	"fmt"
	"log"

	"github.com/miaomiaopu/ipmp-server/internal/config"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/crypto"
	"github.com/miaomiaopu/ipmp-server/internal/router"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化加密模块
	if err := crypto.Init(cfg.Encryption.Key); err != nil {
		log.Fatalf("Failed to init crypto: %v", err)
	}

	// 连接数据库
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	logLevel := logger.Warn
	if cfg.Log.Level == "debug" {
		logLevel = logger.Info
	}
	db.Logger = db.Logger.LogMode(logLevel)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.Project{},
		&model.Task{},
		&model.Requirement{},
		&model.WorkLog{},
		&model.WeeklyReport{},
		&model.UserAIConfig{},
		&model.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// 设置路由
	r := router.Setup(db, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s (mode=%s db=%s)", addr, cfg.Server.Mode, cfg.Database.Type)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	dsn := cfg.Database.DSN()

	switch cfg.Database.Driver() {
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return db, nil
}
