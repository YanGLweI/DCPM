package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ldap-password-manager/config"
	"ldap-password-manager/handler"
	"ldap-password-manager/middleware"
	"ldap-password-manager/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 注册中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 初始化服务
	ldapService := service.NewLDAPService()
	passwordService := service.NewPasswordService(ldapService)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(ldapService)
	passwordHandler := handler.NewPasswordHandler(passwordService)

	// 健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "服务运行正常",
		})
	})

	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关（无需 JWT，但有限流）
		v1.POST("/auth/login", middleware.RateLimit(), authHandler.Login)

		// 需要 JWT 认证的路由
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/info", authHandler.GetUserInfo)
		}

		// 密码修改（密码过期时也可调用，需要特殊处理）
		v1.POST("/password/change", passwordHandler.ChangePassword)
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
