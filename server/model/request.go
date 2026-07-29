package model

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginResponseData 登录成功响应数据
type LoginResponseData struct {
	Status               string `json:"status"`
	Message              string `json:"message"`
	Token                string `json:"token,omitempty"`
	Username             string `json:"username,omitempty"`
	PasswordExpiresAt    string `json:"password_expires_at,omitempty"`
	DaysRemaining        int    `json:"days_remaining,omitempty"`
	PasswordNeverExpires bool   `json:"password_never_expires,omitempty"`
}

// UserInfoResponseData 用户信息响应数据
type UserInfoResponseData struct {
	Username             string `json:"username"`
	PasswordExpiresAt    string `json:"password_expires_at"`
	DaysRemaining        int    `json:"days_remaining"`
	PasswordNeverExpires bool   `json:"password_never_expires"`
}

// AuthResult 认证结果
type AuthResult struct {
	Status               string // ok, expired, error
	PasswordExpiresAt    string
	DaysRemaining        int
	PasswordNeverExpires bool
}
