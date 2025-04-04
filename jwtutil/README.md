# JWTUtil

JWTUtil 是一个轻量级的 Go 语言 JWT 工具包，旨在简化 JWT 的配置和使用流程，加快小型项目的开发速度。

## 特性

- **简单易用**：简洁的 API 接口，只需几行代码即可完成 JWT 的生成和验证
- **灵活配置**：使用函数选项模式，允许灵活配置 JWT 参数
- **合理默认值**：提供默认配置，减少必要的配置项
- **基于官方包**：基于广泛使用的 `github.com/golang-jwt/jwt/v5` 包构建

## 安装

```bash
go get github.com/xcc-d/jwtutil
```

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/your-username/jwtutil"
)

func main() {
    // 初始化 JWT 工具
    jwtUtil := jwtutil.New(
        jwtutil.WithSecret([]byte("your-secret-key")),
    )

    // 创建自定义 Claims
    claims := &jwt.RegisteredClaims{
        Subject:   "user-123",
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
    }

    // 生成 Token
    token, err := jwtUtil.GenerateToken(claims)
    if err != nil {
        fmt.Println("生成 Token 失败:", err)
        return
    }
    fmt.Println("生成的 Token:", token)

    // 解析 Token
    parsedClaims := &jwt.RegisteredClaims{}
    err = jwtUtil.ParseToken(token, parsedClaims)
    if err != nil {
        fmt.Println("解析 Token 失败:", err)
        return
    }
    fmt.Println("Subject:", parsedClaims.Subject)

    // 刷新 Token
    newToken, err := jwtUtil.RefreshToken(token, parsedClaims)
    if err != nil {
        fmt.Println("刷新 Token 失败:", err)
        return
    }
    fmt.Println("刷新后的 Token:", newToken)
}
```

### 自定义 Claims

```go
package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/xcc-d/jwtutil"
)

// UserClaims 自定义 Claims
type UserClaims struct {
    jwt.RegisteredClaims
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
}

func main() {
    // 初始化 JWT 工具，自定义配置
    jwtUtil := jwtutil.New(
        jwtutil.WithSecret([]byte("your-secret-key")),
        jwtutil.WithExpiresIn(24 * time.Hour),
        jwtutil.WithIssuer("my-app"),
    )

    // 创建自定义 Claims
    claims := &UserClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
        UserID:   "123456",
        Username: "johndoe",
        Role:     "admin",
    }

    // 生成 Token
    token, err := jwtUtil.GenerateToken(claims)
    if err != nil {
        fmt.Println("生成 Token 失败:", err)
        return
    }

    // 解析 Token
    parsedClaims := &UserClaims{}
    err = jwtUtil.ParseToken(token, parsedClaims)
    if err != nil {
        fmt.Println("解析 Token 失败:", err)
        return
    }
    fmt.Println("用户名:", parsedClaims.Username)
    fmt.Println("角色:", parsedClaims.Role)
}
```

## 配置选项

JWTUtil 提供了多种配置选项来自定义 JWT 的行为：

| 选项 | 说明 | 默认值 |
|------|------|--------|
| WithSecret | 设置签名密钥（必需） | - |
| WithSigningMethod | 设置签名算法 | HS256 |
| WithExpiresIn | 设置过期时间 | 2小时 |
| WithIssuer | 设置签发者 | 空字符串 |
| WithIssuedAt | 是否设置签发时间 | true |

## 错误处理

```go
import "github.com/xcc-d/jwtutil"

func handleToken(tokenString string) {
    // ...
    if err := jwtUtil.ParseToken(tokenString, claims); err != nil {
        if err == jwtutil.ErrInvalidToken {
            // 处理无效 Token
            return
        }
        // 处理其他错误
    }
    // Token 有效
}
```

