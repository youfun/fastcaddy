```markdown
# FastCaddy + Caddy DNS 插件一键部署教程

本教程介绍如何在生产环境下，使用 FastCaddy 管理 Caddy，并通过 Cloudflare DNS 自动申请 HTTPS 证书。重点讲解如何编译带 Cloudflare DNS 插件的 Caddy。

---

## 1. 安装 Go 环境

请确保你的系统已安装 Go 1.18+  
[Go 官方下载](https://golang.org/dl/)

---

## 2. 编译带 Cloudflare DNS 插件的 Caddy

Caddy 默认不带 Cloudflare DNS 插件，需要用 `xcaddy` 工具自定义编译。

### 安装 xcaddy

```bash
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
export PATH=$PATH:$HOME/go/bin
```

### 编译 Caddy（带 Cloudflare DNS）

```bash
xcaddy build --with github.com/caddy-dns/cloudflare
```

编译完成后，会生成一个 `caddy` 可执行文件。

---

## 3. 替换系统 Caddy

将新编译的 `caddy` 文件覆盖原有 Caddy：

```bash
sudo mv ./caddy /usr/local/bin/caddy
sudo chmod +x /usr/local/bin/caddy
sudo systemctl restart caddy
```

---

## 4. 安装 FastCaddy

```bash
git clone https://github.com/youfun/fastcaddy.git
cd fastcaddy
go build -o fastcaddy ./cmd/fastcaddy
```

---

## 5. 配置 Cloudflare API Token

获取 Cloudflare API Token，并设置环境变量：

```bash
export CADDY_CF_TOKEN="your-cloudflare-token"
```

---

## 6. 初始化 Caddy 配置

首次运行建议先初始化并添加一个域名用于dns挑战：

```bash

 ./fastcaddy setup --local

./fastcaddy add-proxy --from example.com --to localhost:8080

```
进行dns挑战
```bash

export CADDY_CF_TOKEN="your-cloudflare-token"

./fastcaddy setup --cf-token $CADDY_CF_TOKEN
```

---

## 7. 添加反向代理和域名

添加你的域名和后端服务：

```bash
./fastcaddy add-proxy --from example.com --to localhost:8080
./fastcaddy add-wildcard --domain example.com
./fastcaddy add-sub-proxy --domain example.com --subdomain api --ports 8081
```

---

## 8. 检查状态

```bash
./fastcaddy status
```

---

## 常见问题

- **没有域名时无法申请证书**：请先添加域名相关的反向代理。
- **Cloudflare DNS 插件未生效**：请确认 Caddy 已用 xcaddy 编译，并已替换系统 Caddy。
- **API Token 权限不足**：请确保 Cloudflare Token 具备 DNS 编辑权限。

---

## 参考

- [Caddy 官方文档](https://caddyserver.com/docs/)
- [Cloudflare DNS 插件](https://github.com/caddy-dns/cloudflare)
- [FastCaddy 项目](https://github.com/youfun/fastcaddy)
```
