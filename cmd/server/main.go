package main

import (
	"fmt"
	"log"

	"github.com/miaomiaopu/ipmp-server/internal/config"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/crypto"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/version"
	"github.com/miaomiaopu/ipmp-server/internal/router"
	"golang.org/x/crypto/bcrypt"
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

	// 初始化默认数据
	seed(db)

	// 设置路由
	r := router.Setup(db, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("IPMP Server %s (commit=%s built=%s)", version.Version, version.GitCommit, version.BuildTime)
	log.Printf("Listening on %s (mode=%s db=%s)", addr, cfg.Server.Mode, cfg.Database.Type)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seed 初始化默认管理员账户
func seed(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), 12)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return
		}
		db.Create(&model.User{
			Username:     "admin",
			PasswordHash: string(hash),
			DisplayName:  "管理员",
			Role:         "admin",
			Status:       "active",
		})
		log.Println("Default admin user created (admin/admin123)")
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
