# JWTUtil - Go JWT工具包 v1.2

企业级JWT解决方案，支持：
- 多密钥自动轮换
- Token验证缓存
- 多种签名算法(HS256/HS384/HS512/RS256/RS384/RS512)
- 高性能Token处理

## 功能特性

- 支持自定义Claims
- 多种签名算法可选(HS256/HS384/HS512/RS256/RS384/RS512)
- Token自动刷新
- 详细的错误处理
- 灵活的配置选项
- 多密钥自动轮换支持
- Token验证缓存(默认关闭)
- 高性能Token处理

## 安装

```bash
go get github.com/xcc-d/jwtutil
```

## 快速开始

```go
import "github.com/xcc-d/jwtutil"

// 初始化
j := jwtutil.New(
    jwtutil.WithSecret([]byte("your-secret-key")),
    jwtutil.WithExpiresIn(2*time.Hour),
    jwtutil.WithIssuer("your-app"),
)

// 自定义Claims
type MyClaims struct {
    UserID int64 `json:"user_id"`
    jwt.RegisteredClaims
}

// 生成Token
claims := &MyClaims{UserID: 123}
token, err := j.GenerateToken(claims)

// 解析Token
parsedClaims := &MyClaims{}
err = j.ParseToken(token, parsedClaims)

// 刷新Token
newToken, err := j.RefreshToken(token, parsedClaims)
```

## 配置选项说明

所有配置函数都在初始化时通过`New()`函数传入，按需使用：

```go
// 典型配置示例
j := jwtutil.New(
    // 必须设置的密钥
    jwtutil.WithSecret([]byte("your-secret-key")),
    
    // 可选配置
    jwtutil.WithExpiresIn(2*time.Hour),
    jwtutil.WithIssuer("your-app"),
    jwtutil.WithValidateClaims(myValidateFunc),
    
    // v1.2 新功能
    jwtutil.WithMultiSecrets([][]byte{
        []byte("old-secret-1"),
        []byte("old-secret-2"),
    }),
    jwtutil.WithCache(true), // 启用缓存
)
```

### 可用选项列表

| 配置函数 | 说明 | 是否必须 | 默认值 | 版本 |
|----------|------|----------|--------|------|
| WithSecret | 设置签名密钥 | 是 | 无 | v1.0 |
| WithSigningMethod | 设置签名算法 | 否 | HS256 | v1.0 |
| WithExpiresIn | 设置Token默认有效期 | 否 | 2小时 | v1.0 |
| WithIssuer | 设置签发者标识 | 否 | 空 | v1.0 |
| WithIssuedAt | 是否自动设置签发时间 | 否 | true | v1.0 |
| WithValidateClaims | 设置Claims验证回调 | 否 | 无 | v1.0 |
| WithMultiSecrets | 设置多密钥支持 | 否 | 无 | v1.2 |
| WithCache | 启用Token验证缓存 | 否 | false | v1.2 |

## Claims验证示例

```go
j := jwtutil.New(
    jwtutil.WithValidateClaims(func(claims jwt.Claims) error {
        if c, ok := claims.(*MyClaims); ok {
            if c.UserID <= 0 {
                return errors.New("invalid user id")
            }
            if c.Role != "admin" {
                return errors.New("permission denied")
            }
        }
        return nil
    })
)
```

## 错误处理

- `ErrInvalidToken`: Token无效
- `ErrTokenExpired`: Token已过期  
- `ErrInvalidSigningMethod`: 签名算法不匹配
- `ErrMissingKey`: 缺少签名密钥
- `ErrInvalidClaims`: Claims无效
- `ErrTokenMalformed`: Token格式错误

## 性能建议

- 对于高并发场景，建议复用JWTUtil实例
- 避免频繁创建新的Claims对象
