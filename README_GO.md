# FastCaddy Go 版本

这是 fastcaddy Python 库的 Go 重写版本，保持了所有原有功能，同时提供了更好的性能和更丰富的命令行接口。

## 功能特性

- 🚀 **高性能**: 使用 Go 语言重写，性能更佳
- 🛠️ **完整功能**: 保持与 Python 版本 100% 功能兼容
- 🖥️ **命令行工具**: 提供易用的 CLI 界面
- 📚 **编程接口**: 支持作为 Go 库使用
- 🔐 **TLS 支持**: 支持 ACME 和内部证书配置
- 🌐 **路由管理**: 完整的反向代理和通配符域名支持
- 📖 **中文注释**: 关键函数和步骤包含详细中文注释

## 安装

### 从源码编译

```bash
git clone https://github.com/youfun/fastcaddy.git
cd fastcaddy
go build -o fastcaddy ./cmd/fastcaddy
```

### 使用 Go Install

```bash
go install github.com/youfun/fastcaddy/cmd/fastcaddy@latest
```

## 命令行使用

### 设置 Caddy 基础配置

#### 本地开发环境（使用内部证书）
```bash
./fastcaddy setup --local
```

#### 生产环境（使用 Let's Encrypt + Cloudflare）
```bash
export CADDY_CF_TOKEN="your-cloudflare-token"
./fastcaddy setup --cf-token $CADDY_CF_TOKEN
```

#### 安装根证书到系统信任存储
```bash
./fastcaddy setup --local --install-trust
```

### 管理反向代理

#### 添加简单反向代理
```bash
./fastcaddy add-proxy --from api.example.com --to localhost:8080
./fastcaddy add-proxy --from web.example.com --to 127.0.0.1:3000
```

#### 删除反向代理
```bash
./fastcaddy del-proxy --id api.example.com
```

### 通配符子域名支持

#### 添加通配符域名
```bash
./fastcaddy add-wildcard --domain example.com
```

#### 添加子域名反向代理
```bash
# 单端口
./fastcaddy add-sub-proxy --domain example.com --subdomain api --ports 8080

# 多端口
./fastcaddy add-sub-proxy --domain example.com --subdomain web --ports 3000,3001

# 指定目标主机
./fastcaddy add-sub-proxy --domain example.com --subdomain db --ports 5432 --host 192.168.1.10
```

### 查看状态
```bash
./fastcaddy status
```

## 编程接口使用

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/youfun/fastcaddy"
)

func main() {
    // 创建 FastCaddy 客户端
    fc := fastcaddy.New()
    
    // 设置本地开发环境
    err := fc.SetupCaddy("", "srv0", true, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 添加反向代理
    err = fc.AddReverseProxy("api.localhost", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    
    // 添加通配符域名
    err = fc.AddWildcardRoute("localhost")
    if err != nil {
        log.Fatal(err)
    }
    
    // 添加子域名反向代理
    err = fc.AddSubReverseProxy("localhost", "web", []string{"3000"}, "localhost")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Caddy 配置完成!")
}
```

### 高级使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/youfun/fastcaddy"
    "github.com/youfun/fastcaddy/pkg/types"
)

func main() {
    fc := fastcaddy.New()
    
    // 直接操作路由
    route := types.Route{
        ID: "custom-route",
        Match: []types.RouteMatch{
            {Host: []string{"custom.example.com"}},
        },
        Handle: []types.Handler{
            {
                Handler: "reverse_proxy",
                Upstreams: []types.Upstream{
                    {Dial: "backend1:8080"},
                    {Dial: "backend2:8080"},
                },
            },
        },
        Terminal: true,
    }
    
    err := fc.Routes.AddRoute(route)
    if err != nil {
        log.Fatal(err)
    }
    
    // 检查配置状态
    if fc.HasPath("/apps/http/servers") {
        fmt.Println("HTTP 服务器已配置")
    }
    
    // 获取配置
    config, err := fc.GetConfig("/apps/http/servers/srv0")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("服务器配置: %+v\n", config)
}
```

## 项目结构

```
fastcaddy/
├── cmd/
│   └── fastcaddy/          # 命令行工具
├── internal/
│   ├── api/               # Caddy API 交互
│   ├── config/            # 配置管理
│   ├── tls/               # TLS 配置
│   ├── routes/            # 路由管理
│   └── utils/             # 工具函数
├── pkg/
│   └── types/             # 公共类型定义
├── fastcaddy.go           # 主要客户端接口
├── go.mod                 # Go 模块定义
└── go.sum                 # 依赖校验和
```

## 核心组件

### API 客户端 (`internal/api`)
- HTTP 客户端封装
- Caddy REST API 交互
- 配置获取和设置
- 错误处理

### 配置管理 (`internal/config`)
- 嵌套配置操作
- 路径和键值转换
- 配置初始化

### TLS 管理 (`internal/tls`)
- ACME 配置 (Let's Encrypt)
- 内部证书配置
- PKI 信任设置

### 路由管理 (`internal/routes`)
- 反向代理配置
- 通配符域名支持
- 子域名路由

### 工具函数 (`internal/utils`)
- 路径处理
- 环境变量获取
- 验证函数

## 与 Python 版本对比

| 功能 | Python 版本 | Go 版本 | 状态 |
|------|-------------|---------|------|
| Caddy API 交互 | ✅ | ✅ | ✅ 完全兼容 |
| 配置管理 | ✅ | ✅ | ✅ 完全兼容 |
| TLS/SSL 配置 | ✅ | ✅ | ✅ 完全兼容 |
| 反向代理 | ✅ | ✅ | ✅ 完全兼容 |
| 通配符域名 | ✅ | ✅ | ✅ 完全兼容 |
| 命令行工具 | ❌ | ✅ | ✨ 新增功能 |
| 类型安全 | ❌ | ✅ | ✨ 新增功能 |
| 性能 | 🐌 | 🚀 | ✨ 显著提升 |

## 环境变量

- `CADDY_CF_TOKEN`: Cloudflare API 令牌
- `CLOUDFLARE_API_TOKEN`: 备用 Cloudflare API 令牌

## 错误处理

Go 版本提供了更严格的错误处理机制：

```go
fc := fastcaddy.New()

// 所有操作都返回明确的错误信息
err := fc.SetupCaddy("", "srv0", true, nil)
if err != nil {
    // 处理错误
    log.Printf("设置失败: %v", err)
    return
}
```

## 并发安全

Go 版本考虑了并发环境下的安全性，可以在多 goroutine 环境中安全使用。

## 贡献

欢迎贡献代码！请确保：

1. 代码遵循 Go 惯例
2. 添加适当的测试
3. 更新文档
4. 关键函数包含中文注释

## 许可证

与原 Python 版本保持相同的许可证。

## 示例脚本

### 完整的 Web 应用部署

```go
package main

import (
    "log"
    "github.com/youfun/fastcaddy"
)

func main() {
    fc := fastcaddy.New()
    
    // 1. 设置生产环境（假设已设置 CADDY_CF_TOKEN 环境变量）
    err := fc.SetupCaddy("", "srv0", false, nil)
    if err != nil {
        log.Fatalf("设置 Caddy 失败: %v", err)
    }
    
    // 2. 添加主站点
    err = fc.AddReverseProxy("example.com", "localhost:3000")
    if err != nil {
        log.Fatalf("添加主站点失败: %v", err)
    }
    
    // 3. 添加通配符支持
    err = fc.AddWildcardRoute("example.com")
    if err != nil {
        log.Fatalf("添加通配符失败: %v", err)
    }
    
    // 4. 添加 API 子域名
    err = fc.AddSubReverseProxy("example.com", "api", []string{"8080"}, "localhost")
    if err != nil {
        log.Fatalf("添加 API 子域名失败: %v", err)
    }
    
    // 5. 添加管理界面
    err = fc.AddSubReverseProxy("example.com", "admin", []string{"9000"}, "localhost")
    if err != nil {
        log.Fatalf("添加管理界面失败: %v", err)
    }
    
    log.Println("✅ Web 应用部署完成!")
    log.Println("访问地址:")
    log.Println("- 主站点: https://example.com")
    log.Println("- API: https://api.example.com")  
    log.Println("- 管理: https://admin.example.com")
}
```