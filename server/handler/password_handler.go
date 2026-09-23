package handler

import (
	"log"
	"net/http"
	"strings"

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
			
		// 直接匹配业务层错误（400 Bad Request）
		var httpStatus int
		var clientMessage string
			
		// 1. 检查是否是旧密码验证错误
		if strings.Contains(errMsg, "旧密码验证失败") {
			httpStatus = http.StatusBadRequest
			clientMessage = "旧密码验证失败"
		} else if errMsg == "新密码不能与旧密码相同" {
			// 2. 新旧密码相同 -> 400
			httpStatus = http.StatusBadRequest
			clientMessage = errMsg
		} else if strings.HasPrefix(errMsg, "密码长度不能") ||
			strings.HasPrefix(errMsg, "密码不能包含") ||
			strings.HasPrefix(errMsg, "密码必须包含") {
			// 3. 密码复杂度相关错误（前端校验返回的错误） -> 400
			httpStatus = http.StatusBadRequest
			clientMessage = errMsg
		} else if strings.HasPrefix(errMsg, "密码修改失败") ||
			strings.HasPrefix(errMsg, "连接 LDAP") ||
			strings.HasPrefix(errMsg, "搜索用户") ||
			strings.HasPrefix(errMsg, "服务账号绑定") ||
			strings.HasPrefix(errMsg, "修改密码失败") {
			// 4. LDAP 服务错误 -> 500
			httpStatus = http.StatusInternalServerError
			clientMessage = "密码修改失败，请联系管理员"
			log.Printf("[审计] 密码修改失败 [%s] IP:%s 错误：%v", req.Username, c.ClientIP(), err)
		} else {
			// 5. 其他未知错误 -> 500
			httpStatus = http.StatusInternalServerError
			clientMessage = "密码修改失败，请联系管理员"
			log.Printf("[审计] 密码修改失败 [%s] IP:%s 错误：%v", req.Username, c.ClientIP(), err)
		}
			
		c.JSON(httpStatus, model.Response{
			Code:    httpStatus,
			Message: clientMessage,
		})
		return
	}

	log.Printf("[审计] 密码修改成功 [%s] IP:%s", req.Username, c.ClientIP())
	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "密码修改成功",
	})
}
