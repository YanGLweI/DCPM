package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/go-ldap/ldap/v3"
	"log"

	"ldap-password-manager/config"
	"ldap-password-manager/model"
)

// LDAPService LDAP 服务
type LDAPService struct{}

// NewLDAPService 创建 LDAP 服务实例
func NewLDAPService() *LDAPService {
	return &LDAPService{}
}

// NewLDAPConnection 创建 LDAP 连接
func (s *LDAPService) NewLDAPConnection() (*ldap.Conn, error) {
	cfg := config.AppConfig.LDAP

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.Insecure,
	}

	// 加载 CA 证书
	if cfg.CertPath != "" {
		caCert, err := os.ReadFile(cfg.CertPath)
		if err != nil {
			log.Printf("警告: 无法读取证书文件 %s: %v", cfg.CertPath, err)
		} else {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}
	}

	conn, err := ldap.DialURL(cfg.Server, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return nil, fmt.Errorf("连接 LDAP 服务器失败: %w", err)
	}

	return conn, nil
}

// Authenticate 用户认证（判断密码是否过期）
func (s *LDAPService) Authenticate(username, password string) (*model.AuthResult, error) {
	cfg := config.AppConfig.LDAP

	conn, err := s.NewLDAPConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	userDN := fmt.Sprintf("%s@%s", username, cfg.DomainSuffix)
	err = conn.Bind(userDN, password)

	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			errStr := err.Error()
			// AD 密码过期错误码: data 773
			if strings.Contains(errStr, "773") || strings.Contains(errStr, "data 773") {
				return &model.AuthResult{Status: "expired"}, nil
			}
			// 账号被锁定: data 775
			if strings.Contains(errStr, "775") || strings.Contains(errStr, "data 775") {
				return &model.AuthResult{Status: "error"}, fmt.Errorf("账号已被锁定")
			}
		}
		return &model.AuthResult{Status: "error"}, fmt.Errorf("账号或密码错误")
	}

	// 认证成功，查询密码过期时间
	expiresAt, daysRemaining, neverExpires, _ := s.GetPasswordExpiry(username)

	return &model.AuthResult{
		Status:               "ok",
		PasswordExpiresAt:    expiresAt,
		DaysRemaining:        daysRemaining,
		PasswordNeverExpires: neverExpires,
	}, nil
}

// GetPasswordExpiry 查询密码过期时间
// 返回值: 过期时间字符串, 剩余天数, 是否永不过期, error
func (s *LDAPService) GetPasswordExpiry(username string) (string, int, bool, error) {
	cfg := config.AppConfig.LDAP

	conn, err := s.NewLDAPConnection()
	if err != nil {
		return "", 0, false, err
	}
	defer conn.Close()

	// 使用服务账号绑定
	err = conn.Bind(cfg.Username, cfg.Password)
	if err != nil {
		return "", 0, false, fmt.Errorf("服务账号绑定失败: %w", err)
	}

	// 搜索用户（包含 userAccountControl 用于判断密码永不过期）
	searchRequest := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(cfg.UserFilter, username),
		[]string{"pwdLastSet", "sAMAccountName", "userAccountControl"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return "", 0, false, fmt.Errorf("搜索用户失败: %w", err)
	}

	if len(result.Entries) == 0 {
		return "", 0, false, fmt.Errorf("用户不存在")
	}

	entry := result.Entries[0]

	// 检查 userAccountControl 中的 DONT_EXPIRE_PASSWORD 标志 (0x10000 = 65536)
	uacStr := entry.GetAttributeValue("userAccountControl")
	var uac int64
	fmt.Sscanf(uacStr, "%d", &uac)
	neverExpires := (uac & 0x10000) != 0

	if neverExpires {
		log.Printf("[INFO] 用户 %s 设置了密码永不过期 (userAccountControl=%d)", username, uac)
		return "永不过期", -1, true, nil
	}

	// 获取 pwdLastSet (Windows FILETIME 格式)
	pwdLastSetStr := entry.GetAttributeValue("pwdLastSet")
	var pwdLastSetInt int64
	if pwdLastSetStr != "" {
		fmt.Sscanf(pwdLastSetStr, "%d", &pwdLastSetInt)
	} else {
		pwdLastSetBytes := entry.GetRawAttributeValue("pwdLastSet")
		if len(pwdLastSetBytes) == 8 {
			pwdLastSetInt = bytesToInt64(pwdLastSetBytes)
		} else if len(pwdLastSetBytes) > 0 {
			fmt.Sscanf(string(pwdLastSetBytes), "%d", &pwdLastSetInt)
		}
	}

	if pwdLastSetInt == 0 {
		return "未知", 0, false, nil
	}
	pwdLastSet := fileTimeToTime(pwdLastSetInt)
	if pwdLastSet.IsZero() {
		return "未知", 0, false, nil
	}

	// 获取域密码最大使用期限
	maxPwdAge := s.getMaxPwdAge(conn)

	// 计算过期时间
	expiresAt := pwdLastSet.Add(maxPwdAge)
	daysRemaining := int(time.Until(expiresAt).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	return expiresAt.Format("2006-01-02 15:04:05"), daysRemaining, false, nil
}

// getMaxPwdAge 获取域密码最大使用期限
func (s *LDAPService) getMaxPwdAge(conn *ldap.Conn) time.Duration {
	cfg := config.AppConfig.LDAP

	searchRequest := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"maxPwdAge"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return s.defaultMaxPwdAge()
	}

	if len(result.Entries) == 0 {
		return s.defaultMaxPwdAge()
	}

	// 获取 maxPwdAge (可能是字符串形式的整数，也可能是原始字节)
	maxPwdAgeStr := result.Entries[0].GetAttributeValue("maxPwdAge")
	var maxPwdAge int64
	if maxPwdAgeStr != "" {
		fmt.Sscanf(maxPwdAgeStr, "%d", &maxPwdAge)
	} else {
		maxPwdAgeBytes := result.Entries[0].GetRawAttributeValue("maxPwdAge")
		if len(maxPwdAgeBytes) == 8 {
			maxPwdAge = bytesToInt64(maxPwdAgeBytes)
		} else if len(maxPwdAgeBytes) > 0 {
			fmt.Sscanf(string(maxPwdAgeBytes), "%d", &maxPwdAge)
		}
	}

	if maxPwdAge >= 0 {
		return s.defaultMaxPwdAge()
	}

	// maxPwdAge 是负数的 100ns 间隔
	// 转换为 time.Duration (100ns = 1e-7 秒)
	return time.Duration(-maxPwdAge*100) * time.Nanosecond
}

// defaultMaxPwdAge 获取默认密码最大使用期限（从配置读取）
func (s *LDAPService) defaultMaxPwdAge() time.Duration {
	days := config.AppConfig.LDAP.MaxPwdAgeDays
	if days <= 0 {
		days = 30 // 默认 30 天
	}
	return time.Duration(days) * 24 * time.Hour
}

// bytesToInt64 将 8 字节切片转换为 int64（小端序）
func bytesToInt64(b []byte) int64 {
	if len(b) != 8 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(b))
}

// fileTimeToTime 将 Windows FILETIME 转换为 time.Time
func fileTimeToTime(fileTime int64) time.Time {
	if fileTime == 0 {
		return time.Time{}
	}
	// FILETIME 单位是 100ns，先转为秒（避免 int64 溢出）
	// 1 秒 = 10,000,000 个 100ns 单位
	seconds := fileTime / 10000000
	// Windows epoch (1601-01-01) 到 Unix epoch (1970-01-01) 相差 11644473600 秒
	unixSeconds := seconds - 11644473600
	return time.Unix(unixSeconds, 0).Local()
}

// ChangePassword 修改用户密码
func (s *LDAPService) ChangePassword(username, oldPwd, newPwd string) error {
	cfg := config.AppConfig.LDAP

	conn, err := s.NewLDAPConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	// 1. 使用服务账号绑定
	err = conn.Bind(cfg.Username, cfg.Password)
	if err != nil {
		return fmt.Errorf("服务账号绑定失败: %w", err)
	}

	// 2. 搜索用户 DN
	searchRequest := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(cfg.UserFilter, username),
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("搜索用户失败: %w", err)
	}

	if len(result.Entries) == 0 {
		return fmt.Errorf("用户不存在")
	}

	userDN := result.Entries[0].DN

	// 3. 编码新密码 (UTF-16LE + 双引号包裹)
	newPwdEncoded := encodePassword(newPwd)

	// 4. 执行修改
	modifyRequest := ldap.NewModifyRequest(userDN, nil)
	modifyRequest.Replace("unicodePwd", []string{string(newPwdEncoded)})

	err = conn.Modify(modifyRequest)
	if err != nil {
		return fmt.Errorf("密码修改失败: %w", err)
	}

	return nil
}

// VerifyPassword 验证用户密码是否正确
// 注意：密码过期时 Bind 会返回 773 错误，但密码本身是正确的，应视为验证通过
func (s *LDAPService) VerifyPassword(username, password string) error {
	cfg := config.AppConfig.LDAP

	conn, err := s.NewLDAPConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDN := fmt.Sprintf("%s@%s", username, cfg.DomainSuffix)
	err = conn.Bind(userDN, password)
	if err != nil {
		// AD 密码过期错误码: data 773，密码正确但已过期，视为验证通过
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			errStr := err.Error()
			if strings.Contains(errStr, "773") || strings.Contains(errStr, "data 773") {
				return nil
			}
		}
		return err
	}
	return nil
}

// encodePassword 将密码编码为 AD 要求的 UTF-16LE 格式
func encodePassword(password string) []byte {
	// AD 要求密码用双引号包裹，并编码为 UTF-16LE
	quoted := "\"" + password + "\""
	utf16Bytes := make([]byte, len(quoted)*2)
	for i, r := range utf16.Encode([]rune(quoted)) {
		binary.LittleEndian.PutUint16(utf16Bytes[i*2:], r)
	}
	return utf16Bytes
}
