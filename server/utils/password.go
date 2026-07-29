package utils

import (
	"fmt"
	"strings"
	"unicode"

	"ldap-password-manager/config"
)

// ValidatePasswordComplexity 校验密码复杂度
func ValidatePasswordComplexity(password, username string) error {
	policy := config.AppConfig.PasswordPolicy

	// 检查最小长度
	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度不能少于 %d 位", policy.MinLength)
	}

	// 不能包含用户名
	if username != "" && strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		return fmt.Errorf("密码不能包含用户名")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	// 统计满足的条件数
	categoryCount := 0
	if hasUpper {
		categoryCount++
	}
	if hasLower {
		categoryCount++
	}
	if hasDigit {
		categoryCount++
	}
	if hasSpecial {
		categoryCount++
	}

	// 至少需要满足 3 类
	if categoryCount < 3 {
		return fmt.Errorf("密码必须包含大写字母、小写字母、数字、特殊字符中的至少 3 类")
	}

	// 根据配置检查必须包含的字符类型
	if policy.RequireUppercase && !hasUpper {
		return fmt.Errorf("密码必须包含大写字母")
	}
	if policy.RequireLowercase && !hasLower {
		return fmt.Errorf("密码必须包含小写字母")
	}
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("密码必须包含数字")
	}
	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("密码必须包含特殊字符")
	}

	return nil
}

// ValidatePasswordDifferent 校验新旧密码不同
func ValidatePasswordDifferent(oldPwd, newPwd string) error {
	if oldPwd == newPwd {
		return fmt.Errorf("新密码不能与旧密码相同")
	}
	return nil
}

// GetPasswordStrength 获取密码强度等级 (1-4)
func GetPasswordStrength(password string) int {
	strength := 0
	if len(password) >= 14 {
		strength++
	}
	if strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		strength++
	}
	if strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		strength++
	}
	if strings.ContainsAny(password, "0123456789") {
		strength++
	}
	if strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
		strength++
	}

	switch {
	case strength <= 2:
		return 1 // 弱
	case strength == 3:
		return 2 // 中
	case strength == 4:
		return 3 // 强
	default:
		return 4 // 很强
	}
}
