package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ldap-password-manager/model"
	"ldap-password-manager/service"
	"ldap-password-manager/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	ldapService *service.LDAPService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(ldapService *service.LDAPService) *AuthHandler {
	return &AuthHandler{
		ldapService: ldapService,
	}
}

// Login 登录验证
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	// 执行认证
	result, err := h.ldapService.Authenticate(req.Username, req.Password)
	if err != nil {
		if result != nil && result.Status == "error" {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code:    401,
				Message: err.Error(),
				Data: model.LoginResponseData{
					Status:  "error",
					Message: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "认证服务异常",
		})
		return
	}

	switch result.Status {
	case "expired":
		// 密码过期
		c.JSON(http.StatusOK, model.Response{
			Code:    200,
			Message: "密码已过期，请修改密码",
			Data: model.LoginResponseData{
				Status:  "expired",
				Message: "密码已过期，请修改密码",
			},
		})

	case "ok":
		// 认证成功，签发 JWT
		token, err := utils.GenerateToken(req.Username)
		if err != nil {
			log.Printf("生成 JWT 失败: %v", err)
			c.JSON(http.StatusInternalServerError, model.Response{
				Code:    500,
				Message: "生成认证令牌失败",
			})
			return
		}

		c.JSON(http.StatusOK, model.Response{
			Code:    200,
			Message: "登录成功",
			Data: model.LoginResponseData{
				Status:               "ok",
				Message:              "登录成功",
				Token:                token,
				Username:             req.Username,
				PasswordExpiresAt:    result.PasswordExpiresAt,
				DaysRemaining:        result.DaysRemaining,
				PasswordNeverExpires: result.PasswordNeverExpires,
			},
		})

	default:
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    401,
			Message: "账号或密码错误",
			Data: model.LoginResponseData{
				Status:  "error",
				Message: "账号或密码错误",
			},
		})
	}
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	usernameStr := username.(string)

	// 查询密码过期时间
	expiresAt, daysRemaining, neverExpires, err := h.ldapService.GetPasswordExpiry(usernameStr)
	if err != nil {
		log.Printf("查询密码过期时间失败: %v", err)
		c.JSON(http.StatusOK, model.Response{
			Code:    200,
			Message: "获取用户信息成功",
			Data: model.UserInfoResponseData{
				Username:             usernameStr,
				PasswordExpiresAt:    "未知",
				DaysRemaining:        0,
				PasswordNeverExpires: false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "获取用户信息成功",
		Data: model.UserInfoResponseData{
			Username:             usernameStr,
			PasswordExpiresAt:    expiresAt,
			DaysRemaining:        daysRemaining,
			PasswordNeverExpires: neverExpires,
		},
	})
}
