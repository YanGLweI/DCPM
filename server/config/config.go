package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig   `mapstructure:"server"`
	JWT            JWTConfig      `mapstructure:"jwt"`
	LDAP           LDAPConfig     `mapstructure:"ldap"`
	PasswordPolicy PasswordPolicy `mapstructure:"password_policy"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expiry time.Duration `mapstructure:"expiry"`
}

type LDAPConfig struct {
	Server         string `mapstructure:"server"`
	BaseDN         string `mapstructure:"base_dn"`
	DomainSuffix   string `mapstructure:"domain_suffix"`
	UseTLS         bool   `mapstructure:"use_tls"`
	Insecure       bool   `mapstructure:"insecure"`
	UserFilter     string `mapstructure:"user_filter"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	CertPath       string `mapstructure:"cert_path"`
	MaxPwdAgeDays  int    `mapstructure:"max_pwd_age_days"` // 密码最大使用天数（默认值，当无法从域控获取时使用）
}

type PasswordPolicy struct {
	MinLength        int  `mapstructure:"min_length"`
	RequireUppercase bool `mapstructure:"require_uppercase"`
	RequireLowercase bool `mapstructure:"require_lowercase"`
	RequireDigit     bool `mapstructure:"require_digit"`
	RequireSpecial   bool `mapstructure:"require_special"`
}

var AppConfig *Config

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 支持通过环境变量覆盖敏感配置（优先级高于配置文件）
	if v := os.Getenv("LDAP_SERVER"); v != "" {
		cfg.LDAP.Server = v
	}
	if v := os.Getenv("LDAP_USERNAME"); v != "" {
		cfg.LDAP.Username = v
	}
	if v := os.Getenv("LDAP_PASSWORD"); v != "" {
		cfg.LDAP.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}

	AppConfig = &cfg
	return &cfg, nil
}
