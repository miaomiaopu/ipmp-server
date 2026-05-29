package router

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/config"
	"github.com/miaomiaopu/ipmp-server/internal/handler"
	"github.com/miaomiaopu/ipmp-server/internal/middleware"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/jwt"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"github.com/miaomiaopu/ipmp-server/internal/service"
	"gorm.io/gorm"
)

// Setup 初始化 Gin 路由，注册所有中间件、API 端点
// 分层组装: repository → service → handler → route
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// 安全基础中间件
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.MethodAllow())
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins()))
	r.Use(middleware.GlobalRateLimit())
	r.Use(middleware.MaxBodySize(10 << 20)) // 10MB
	r.Use(middleware.SensitiveLog())

	// JWT 管理器
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshExpire,
	)

	// 初始化层
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	authSvc := service.NewAuthService(userRepo, jwtManager)
	customerSvc := service.NewCustomerService(customerRepo)
	projectSvc := service.NewProjectService(projectRepo)

	authH := handler.NewAuthHandler(authSvc)
	customerH := handler.NewCustomerHandler(customerSvc)
	projectH := handler.NewProjectHandler(projectSvc)

	api := r.Group("/api/v1")
	{
		// 认证（登录限流独立 + 审计）
		auth := api.Group("/auth")
		auth.Use(middleware.AuditLogger(db))
		{
			auth.POST("/login", middleware.LoginRateLimit(), authH.Login)
			auth.POST("/refresh", authH.Refresh)
			auth.POST("/logout", middleware.AuthRequired(jwtManager), authH.Logout)
			auth.GET("/me", middleware.AuthRequired(jwtManager), authH.Me)
		}

		// 客户（认证 + 审计，admin 不可操作）
		customers := api.Group("/customers",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("manager", "user"),
			middleware.AuditLogger(db))
		{
			customers.GET("", customerH.List)
			customers.GET("/:id", customerH.GetByID)
			customers.POST("", customerH.Create)
			customers.POST("/:id/update", customerH.Update)
			customers.POST("/:id/delete", customerH.Delete)
		}

		// 项目（认证 + 审计，admin 不可操作）
		projects := api.Group("/projects",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("manager", "user"),
			middleware.AuditLogger(db))
		{
			projects.GET("", projectH.List)
			projects.GET("/:id", projectH.GetByID)
			projects.POST("", projectH.Create)
			projects.POST("/:id/update", projectH.Update)
			projects.POST("/:id/delete", projectH.Delete)
		}
	}

	return r
}
