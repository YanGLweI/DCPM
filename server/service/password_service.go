package service

import (
	"fmt"

	"ldap-password-manager/utils"
)

// PasswordService 密码业务逻辑服务
type PasswordService struct {
	ldapService *LDAPService
}

// NewPasswordService 创建密码服务实例
func NewPasswordService(ldapService *LDAPService) *PasswordService {
	return &PasswordService{
		ldapService: ldapService,
	}
}

// ChangePassword 修改密码（包含业务校验）
func (s *PasswordService) ChangePassword(username, oldPwd, newPwd string) error {
	// 1. 校验新旧密码不同
	if err := utils.ValidatePasswordDifferent(oldPwd, newPwd); err != nil {
		return err
	}

	// 2. 校验新密码复杂度
	if err := utils.ValidatePasswordComplexity(newPwd, username); err != nil {
		return err
	}

	// 3. 验证旧密码（如果密码未过期）
	if err := s.ldapService.VerifyPassword(username, oldPwd); err != nil {
		return fmt.Errorf("旧密码验证失败")
	}

	// 4. 修改密码
	if err := s.ldapService.ChangePassword(username, oldPwd, newPwd); err != nil {
		return err
	}

	return nil
}
