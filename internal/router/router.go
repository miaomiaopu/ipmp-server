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
taskRepo := repository.NewTaskRepository(db)
	requirementRepo := repository.NewRequirementRepository(db)
	workLogRepo := repository.NewWorkLogRepository(db)
	aiConfigRepo := repository.NewUserAIConfigRepository(db)

	authSvc := service.NewAuthService(userRepo, jwtManager)
	userSvc := service.NewUserService(userRepo)
	customerSvc := service.NewCustomerService(customerRepo)
	projectSvc := service.NewProjectService(projectRepo)
taskSvc := service.NewTaskService(taskRepo)
	requirementSvc := service.NewRequirementService(requirementRepo)
	workLogSvc := service.NewWorkLogService(workLogRepo)
	aiConfigSvc := service.NewUserAIConfigService(aiConfigRepo)

	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	customerH := handler.NewCustomerHandler(customerSvc)
	projectH := handler.NewProjectHandler(projectSvc)
taskH := handler.NewTaskHandler(taskSvc)
	requirementH := handler.NewRequirementHandler(requirementSvc)
	workLogH := handler.NewWorkLogHandler(workLogSvc)
	aiConfigH := handler.NewUserAIConfigHandler(aiConfigSvc)

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
			auth.POST("/change-password", middleware.AuthRequired(jwtManager), userH.ChangePassword)
		}

		// 用户管理（仅 admin，含审计）
		users := api.Group("/users",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("admin"),
			middleware.AuditLogger(db))
		{
			users.GET("", userH.List)
			users.GET("/:id", userH.GetByID)
			users.POST("", userH.Create)
			users.POST("/:id/update", userH.Update)
			users.POST("/:id/delete", userH.Delete)
			users.POST("/:id/force-delete", userH.ForceDelete)
			users.POST("/:id/restore", userH.Restore)
			users.POST("/:id/reset-password", userH.ResetPassword)
		}

		// 客户（认证 + 审计，admin 不可操作）
		customers := api.Group("/customers",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("user"),
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
			middleware.RequireRole("user"),
			middleware.AuditLogger(db))
		{
			projects.GET("", projectH.List)
			projects.GET("/:id", projectH.GetByID)
			projects.POST("", projectH.Create)
			projects.POST("/:id/update", projectH.Update)
			projects.POST("/:id/delete", projectH.Delete)
		}

		// 任务（认证 + 审计，admin 不可操作）
		tasks := api.Group("/tasks",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("user"),
			middleware.AuditLogger(db))
		{
			tasks.GET("", taskH.List)
			tasks.GET("/:id", taskH.GetByID)
			tasks.POST("", taskH.Create)
			tasks.POST("/:id/update", taskH.Update)
			tasks.POST("/:id/status", taskH.UpdateStatus)
			tasks.POST("/:id/delete", taskH.Delete)
		}

		// 需求（认证 + 审计，admin 不可操作）
		reqs := api.Group("/requirements",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("user"),
			middleware.AuditLogger(db))
		{
			reqs.GET("", requirementH.List)
			reqs.GET("/:id", requirementH.GetByID)
			reqs.POST("", requirementH.Create)
			reqs.POST("/:id/update", requirementH.Update)
			reqs.POST("/:id/delete", requirementH.Delete)
		}

		// 工时（认证 + 审计，admin 不可操作）
		wl := api.Group("/work-logs",
			middleware.AuthRequired(jwtManager),
			middleware.RequireRole("user"),
			middleware.AuditLogger(db))
		{
			wl.GET("", workLogH.List)
			wl.GET("/stats", workLogH.Stats)
			wl.GET("/:id", workLogH.GetByID)
			wl.POST("", workLogH.Create)
			wl.POST("/:id/update", workLogH.Update)
			wl.POST("/:id/delete", workLogH.Delete)
		}

		// AI 配置（认证，所有用户可配置自己的 Key）
		aiConfig := api.Group("/ai-config",
			middleware.AuthRequired(jwtManager),
			middleware.AuditLogger(db))
		{
			aiConfig.GET("", aiConfigH.Get)
			aiConfig.POST("/update", aiConfigH.Update)
			aiConfig.POST("/delete", aiConfigH.Delete)
		}
	}

	return r
}
