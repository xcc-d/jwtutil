package jwtutil

import (
	"sync"
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

var (
	tokenCache = &sync.Map{}
)

// ParseToken 解析并验证Token
func (j *JWTUtil) ParseToken(tokenString string, claims jwt.Claims) error {
	// 检查缓存
	if j.options.enableCache {
		if cached, ok := tokenCache.Load(tokenString); ok {
			if err, ok := cached.(error); ok {
				return err
			}
			return nil
		}
	}

	// 多密钥验证
	verifyFn := func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.options.signingMethod {
			return nil, ErrInvalidSigningMethod
		}

		// 尝试所有密钥
		secrets := append([][]byte{j.options.secret}, j.options.secrets...)
		var lastErr error
		for _, secret := range secrets {
			if len(secret) == 0 {
				continue
			}
			// 直接使用签名字符串验证
			signingString := token.Raw[:len(token.Raw)-len(token.Signature)-1]
			if err := token.Method.Verify(signingString, token.Signature, secret); err == nil {
				return secret, nil
			} else {
				lastErr = err
			}
		}
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, ErrMissingKey
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, verifyFn)

	// 缓存结果
	if j.options.enableCache {
		if err == nil {
			tokenCache.Store(tokenString, true)
		} else {
			tokenCache.Store(tokenString, err)
		}
	}

	if err != nil {
		switch err {
		case jwt.ErrTokenMalformed:
			return ErrTokenMalformed
		case jwt.ErrTokenExpired:
			return ErrTokenExpired
		case jwt.ErrTokenNotValidYet:
			return ErrInvalidToken
		default:
			return err
		}
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

	// 重置签发时间和过期时间 1
	if regClaims, ok := claims.(*jwt.RegisteredClaims); ok {
		regClaims.IssuedAt = jwt.NewNumericDate(time.Now())
		if j.options.expiresIn > 0 {
			regClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.options.expiresIn))
		}
	}

	return j.GenerateToken(claims)
}
