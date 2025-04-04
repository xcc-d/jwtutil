package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil JWT工具封装
type JWTUtil struct {
	options *Options
}

// New 创建JWTUtil实例
func New(opts ...Option) *JWTUtil {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &JWTUtil{options: options}
}

// GenerateToken 生成JWT Token
func (j *JWTUtil) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(j.options.signingMethod, claims)

	// 设置签发时间
	if j.options.issuedAt {
		token.Claims.(jwt.Claims).(*jwt.RegisteredClaims).IssuedAt = jwt.NewNumericDate(time.Now())
	}

	// 设置签发者
	if j.options.issuer != "" {
		token.Claims.(jwt.Claims).(*jwt.RegisteredClaims).Issuer = j.options.issuer
	}

	return token.SignedString(j.options.secret)
}

// ParseToken 解析并验证Token
func (j *JWTUtil) ParseToken(tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.options.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

// RefreshToken 刷新Token
func (j *JWTUtil) RefreshToken(tokenString string, claims jwt.Claims) (string, error) {
	if err := j.ParseToken(tokenString, claims); err != nil {
		return "", err
	}

	// 重置签发时间和过期时间
	if regClaims, ok := claims.(*jwt.RegisteredClaims); ok {
		regClaims.IssuedAt = jwt.NewNumericDate(time.Now())
		if j.options.expiresIn > 0 {
			regClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.options.expiresIn))
		}
	}

	return j.GenerateToken(claims)
}
