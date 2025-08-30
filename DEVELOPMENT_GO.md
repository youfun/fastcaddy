# fastcaddy - Golang 重构开发文档

## 项目概述

fastcaddy 是一个用于简化 Caddy Web 服务器 API 使用的 Python 库。它提供了一组简洁的函数来配置和管理 Caddy 服务器，特别是用于设置反向代理、SSL 证书和路由规则。

本项目旨在将 fastcaddy 重构为 Golang 版本，以提供更好的性能和更广泛的适用性。

## 核心功能分析

### 1. Caddy API 交互
- 与 Caddy 服务器的 REST API 进行交互
- 通过 HTTP 请求管理配置和路由

### 2. 配置管理
- 获取和设置 Caddy 配置路径
- 管理 ID 和路径的映射关系

### 3. SSL/TLS 配置
- 支持 ACME 协议（用于生产环境）
- 支持内部证书（用于本地开发）

### 4. 路由管理
- 添加和删除路由规则
- 配置反向代理
- 支持通配符子域名

### 5. 服务部署
- 提供 systemd 服务配置
- 自动化部署脚本

## 核心函数列表

### API 交互函数
1. `get_id(path)` - 生成 ID 的完整 URL
2. `get_path(path)` - 生成配置的完整 URL
3. `gid(path)` - 获取指定路径的 ID 配置
4. `gcfg(path, method)` - 获取指定路径的配置
5. `pid(d, path, method)` - 将配置放入指定 ID 路径
6. `pcfg(d, path, method)` - 将配置放入指定配置路径
7. `has_id(id)` - 检查 ID 是否已设置
8. `has_path(path)` - 检查路径是否已设置

### 配置辅助函数
9. `nested_setdict(sd, value, *keys)` - 在嵌套字典中设置值
10. `path2keys(path)` - 将路径分割为键列表
11. `keys2path(*keys)` - 将键列表连接为路径
12. `nested_setcfg(value, *keys)` - 在配置中设置嵌套值
13. `init_path(path, skip)` - 初始化配置路径

### SSL/TLS 配置函数
14. `get_acme_config(token)` - 获取 ACME 配置
15. `add_tls_internal_config()` - 添加内部 TLS 配置
16. `add_acme_config(cf_token)` - 添加 ACME 配置
17. `setup_pki_trust(install_trust)` - 配置 PKI 证书信任

### 路由管理函数
18. `init_routes(srv_name, skip)` - 初始化 HTTP 路由配置
19. `setup_caddy(cf_token, srv_name, local, install_trust)` - 设置 Caddy 基本配置
20. `add_route(route)` - 添加路由规则
21. `del_id(id)` - 删除指定 ID 的路由
22. `add_reverse_proxy(from_host, to_url)` - 添加反向代理路由
23. `add_wildcard_route(domain)` - 添加通配符子域名路由
24. `add_sub_reverse_proxy(domain, subdomain, port, host)` - 添加子域名反向代理

## 主要变量和常量

### 路径常量
- `automation_path = "/apps/tls/automation"` - TLS 自动化配置路径
- `srvs_path = "/apps/http/servers"` - HTTP 服务器配置路径
- `rts_path = "/apps/http/servers/srv0/routes"` - 路由配置路径

### 环境变量
- `CLOUDFLARE_API_TOKEN` - Cloudflare API 令牌
- `CADDY_CF_TOKEN` - Caddy 使用的 Cloudflare 令牌

## Golang 重构设计

### 包结构设计
```
fastcaddy/
├── cmd/                 # 命令行工具
│   └── fastcaddy/
├── internal/
│   ├── api/            # Caddy API 交互
│   ├── config/         # 配置管理
│   ├── tls/            # TLS 配置
│   ├── routes/         # 路由管理
│   └── utils/          # 工具函数
├── pkg/
│   └── types/          # 公共类型定义
├── docs/               # 文档
├── examples/           # 示例代码
├── go.mod
└── go.sum
```

### 核心模块设计

#### 1. API 模块 (`internal/api`)
- 实现与 Caddy REST API 的交互
- 封装 HTTP 客户端
- 提供配置获取和设置功能

#### 2. 配置模块 (`internal/config`)
- 管理配置结构体
- 实现路径和 ID 的映射
- 提供配置验证功能

#### 3. TLS 模块 (`internal/tls`)
- 实现 ACME 和内部证书配置
- 管理证书颁发机构信任设置

#### 4. 路由模块 (`internal/routes`)
- 管理路由规则
- 实现反向代理配置
- 支持通配符子域名

#### 5. 工具模块 (`internal/utils`)
- 实现辅助函数
- 提供路径处理和嵌套字典操作

### 主要类型定义

```go
// Caddy 配置结构
type CaddyConfig struct {
    Apps map[string]interface{} `json:"apps"`
}

// 路由规则结构
type Route struct {
    ID        string        `json:"@id,omitempty"`
    Match     []RouteMatch  `json:"match"`
    Handle    []Handler     `json:"handle"`
    Terminal  bool          `json:"terminal"`
}

// 路由匹配规则
type RouteMatch struct {
    Host []string `json:"host"`
}

// 处理器结构
type Handler struct {
    Handler   string     `json:"handler"`
    Upstreams []Upstream `json:"upstreams,omitempty"`
    Routes    []Route    `json:"routes,omitempty"`
}

// 上游服务器
type Upstream struct {
    Dial string `json:"dial"`
}
```

### 命令行接口设计

```
fastcaddy [command] [flags]

Commands:
  setup        设置 Caddy 基本配置
  add-proxy    添加反向代理
  del-proxy    删除反向代理
  add-wildcard 添加通配符子域名
  status       查看 Caddy 配置状态
  help         显示帮助信息

Flags:
  -h, --help   显示帮助信息
```

### 使用示例

#### 设置 Caddy
```bash
# 为生产环境设置 Caddy（使用 Cloudflare）
fastcaddy setup --cf-token $CADDY_CF_TOKEN

# 为本地开发设置 Caddy（使用内部证书）
fastcaddy setup --local
```

#### 添加反向代理
```bash
# 添加单个反向代理
fastcaddy add-proxy --from example.com --to localhost:8080

# 添加通配符子域名反向代理
fastcaddy add-wildcard --domain example.com
fastcaddy add-proxy --from sub.example.com --to localhost:8080
```

## 部署方案

### systemd 服务配置
提供与原项目相同的 systemd 服务配置脚本，确保与现有部署流程兼容。

### Docker 支持
考虑提供 Docker 镜像支持，简化部署流程。

## 开发计划

### 第一阶段：核心功能实现
1. 实现 Caddy API 交互模块
2. 实现配置管理模块
3. 实现基本的路由管理功能

### 第二阶段：TLS 和安全功能
1. 实现 ACME 配置功能
2. 实现内部证书配置
3. 实现 PKI 信任配置

### 第三阶段：高级路由功能
1. 实现反向代理功能
2. 实现通配符子域名支持
3. 实现多端口支持

### 第四阶段：工具和部署
1. 实现命令行工具
2. 提供部署脚本
3. 编写完整文档和示例

## 注意事项

1. **错误处理**：Golang 版本需要更严格的错误处理机制
2. **并发安全**：考虑在多 goroutine 环境下的安全性
3. **配置验证**：在发送到 Caddy API 前验证配置的有效性
4. **向后兼容**：保持与 Python 版本相似的 API 设计，便于迁移


添加反向代理的详细流程

  主要流程

  添加反向代理的核心函数是 add_reverse_proxy(from_host, to_url)，但这个函数依赖于多个其他函数协同工作。

  涉及的函数及其作用

  1. add_reverse_proxy(from_host, to_url)
  作用: 创建一个反向代理处理器
  流程:
   1. 检查是否存在同名的主机ID，如果存在则删除
   2. 构造路由配置对象
   3. 调用 add_route(route) 添加路由

    1 def add_reverse_proxy(from_host, to_url):
    2     "Create a reverse proxy handler"
    3     if has_id(from_host): del_id(from_host)
    4     route = {
    5         "handle": [{ "handler": "reverse_proxy",
    6                     "upstreams": [{"dial": to_url}] }],
    7         "match": [{"host": [from_host]}],
    8         "@id": from_host,
    9         "terminal": True
   10     }
   11     add_route(route)

  2. has_id(id)
  作用: 检查指定ID是否已经设置
  依赖函数: gid(id)
  流程:
   1. 调用 gid(id) 尝试获取指定ID的配置
   2. 如果成功获取，返回True
   3. 如果发生HTTP错误（如404），返回False

   1 def has_id(id):
   2     "Check if `id` is set up"
   3     try: gid(id)
   4     except HTTPStatusError: return False
   5     return True

  3. gid(path='/')
  作用: 获取指定路径的ID配置
  依赖函数: get_id(path), xget (httpx.get)
  流程:
   1. 调用 get_id(path) 构造完整的ID URL
   2. 发送HTTP GET请求获取配置
   3. 将JSON响应转换为对象并返回

   1 def gid(path='/'):
   2     "Gets the id at `path`"
   3     response = xget(get_id(path))
   4     response.raise_for_status()
   5     return dict2obj(response.json())

  4. get_id(path)
  作用: 生成ID的完整URL
  流程:
   1. 确保路径以 '/' 开头和结尾
   2. 返回格式为 http://localhost:2019/id/{path} 的URL

   1 def get_id(path):
   2     "Get a ID full URL from a path"
   3     if path[0 ]!='/': path = '/'+path
   4     if path[-1]!='/': path = path+'/'
   5     return f'http://localhost:2019/id{path}'

  5. del_id(id)
  作用: 删除指定ID的路由
  依赖函数: get_id(id), xdelete (httpx.delete)
  流程:
   1. 调用 get_id(id) 获取完整ID URL
   2. 发送HTTP DELETE请求删除配置

   1 def del_id(id):
   2     "Delete route for `id` (e.g. a host)"
   3     xdelete(get_id(id))

  6. add_route(route)
  作用: 将路由配置添加到Caddy配置中
  依赖函数: pcfg(d, path, method)
  流程:
   1. 调用 pcfg(route, rts_path) 将路由配置发送到路由配置路径
   2. rts_path 定义为 /apps/http/servers/srv0/routes

   1 def add_route(route):
   2     "Add `route` dict to config"
   3     return pcfg(route, rts_path)

  7. pcfg(d, path='/', method='post')
  作用: 将配置对象放入指定的配置路径
  依赖函数: get_path(path), obj2dict(d) (来自fastcore), HTTP方法函数
  流程:
   1. 根据method参数获取对应的HTTP方法函数（post, patch等）
   2. 调用 get_path(path) 获取完整配置URL
   3. 将配置对象转换为字典
   4. 发送HTTP请求将配置写入Caddy

   1 def pcfg(d, path='/', method='post'):
   2     "Puts the config `d` into `path`"
   3     f = getattr(httpx, method)
   4     response = f(get_path(path), json=obj2dict(d))
   5     try: response.raise_for_status()
   6     except Exception as e:
   7         e.add_note(f"Error: '{json.loads(response.text)['error']}'")
   8         raise
   9     return response.text or None

  8. get_path(path)
  作用: 生成配置的完整URL
  流程:
   1. 确保路径以 '/' 开头和结尾
   2. 返回格式为 http://localhost:2019/config/{path} 的URL

   1 def get_path(path):
   2     "Get a config full URL from a path"
   3     if path[0 ]!='/': path = '/'+path
   4     if path[-1]!='/': path = path+'/'
   5     return f'http://localhost:2019/config{path}'

  完整流程图

    1 add_reverse_proxy(from_host, to_url)
    2 │
    3 ├─► has_id(from_host)
    4 │   │
    5 │   └─► gid(from_host)
    6 │       │
    7 │       └─► get_id(from_host)
    8 │           │
    9 │           └─► 构造URL: http://localhost:2019/id/{from_host}/
   10 │
   11 ├─► (如果ID存在) del_id(from_host)
   12 │   │
   13 │   └─► get_id(from_host)
   14 │       │
   15 │       └─► 发送DELETE请求删除现有配置
   16 │
   17 ├─► 构造路由配置对象
   18 │   {
   19 │     "handle": [{ "handler": "reverse_proxy", "upstreams": [{"dial": to_url}] }],
   20 │     "match": [{"host": [from_host]}],
   21 │     "@id": from_host,
   22 │     "terminal": True
   23 │   }
   24 │
   25 └─► add_route(route)
   26     │
   27     └─► pcfg(route, rts_path)
   28         │
   29         ├─► get_path(rts_path)
   30         │   │
   31         │   └─► 构造URL: http://localhost:2019/config/apps/http/servers/srv0/routes/
   32         │
   33         └─► 发送POST请求将路由配置写入Caddy

  关键变量

  路径常量
   - srvs_path = '/apps/http/servers' - HTTP服务器配置路径
   - rts_path = srvs_path+'/srv0/routes' - 路由配置路径（即 /apps/http/servers/srv0/routes）

  错误处理

  在整个流程中，错误处理主要通过以下方式实现：
   1. response.raise_for_status() - 检查HTTP响应状态
   2. try...except HTTPStatusError - 捕获HTTP状态错误（如404）
   3. try...except Exception - 捕获其他异常并在错误信息中添加详细说明
   