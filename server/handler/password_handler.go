package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ldap-password-manager/model"
	"ldap-password-manager/service"
)

// PasswordHandler 密码处理器
type PasswordHandler struct {
	passwordService *service.PasswordService
}

// NewPasswordHandler 创建密码处理器
func NewPasswordHandler(passwordService *service.PasswordService) *PasswordHandler {
	return &PasswordHandler{
		passwordService: passwordService,
	}
}

// ChangePassword 修改密码
func (h *PasswordHandler) ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	// 执行密码修改
	err := h.passwordService.ChangePassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		errMsg := err.Error()

		// 根据错误类型返回不同的状态码
		if errMsg == "新密码不能与旧密码相同" ||
			errMsg == "旧密码验证失败" {
			c.JSON(http.StatusBadRequest, model.Response{
				Code:    400,
				Message: errMsg,
			})
			return
		}

		// 密码复杂度校验失败
		if len(errMsg) > 4 && (errMsg[:4] == "密码" || errMsg[:2] == "新密") {
			c.JSON(http.StatusBadRequest, model.Response{
				Code:    400,
				Message: errMsg,
			})
			return
		}

		log.Printf("[审计] 密码修改失败 [%s] IP:%s 错误: %v", req.Username, c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "密码修改失败，请联系管理员",
		})
		return
	}

	log.Printf("[审计] 密码修改成功 [%s] IP:%s", req.Username, c.ClientIP())
	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "密码修改成功",
	})
}
