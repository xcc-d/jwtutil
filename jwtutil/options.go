package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Options JWT配置选项
type Options struct {
	secret         []byte
	signingMethod  jwt.SigningMethod
	expiresIn      time.Duration
	issuer         string
	issuedAt       bool
	validateClaims func(claims jwt.Claims) error // Claims验证回调
}

// Option 配置函数类型
type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		signingMethod: jwt.SigningMethodHS256,
		expiresIn:     2 * time.Hour,
		issuedAt:      true,
	}
}

// WithSecret 设置密钥
func WithSecret(secret []byte) Option {
	return func(o *Options) {
		o.secret = secret
	}
}

// WithSigningMethod 设置签名方法
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *Options) {
		o.signingMethod = method
	}
}

// WithExpiresIn 设置过期时间
func WithExpiresIn(d time.Duration) Option {
	return func(o *Options) {
		o.expiresIn = d
	}
}

// WithIssuer 设置签发者
func WithIssuer(issuer string) Option {
	return func(o *Options) {
		o.issuer = issuer
	}
}

// WithIssuedAt 是否设置签发时间
func WithIssuedAt(enable bool) Option {
	return func(o *Options) {
		o.issuedAt = enable
	}
}

// WithValidateClaims 设置Claims验证回调
func WithValidateClaims(validator func(claims jwt.Claims) error) Option {
	return func(o *Options) {
		o.validateClaims = validator
	}
}
