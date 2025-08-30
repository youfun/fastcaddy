package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youfun/fastcaddy"
	"github.com/youfun/fastcaddy/internal/utils"
)

var (
	// 全局变量存储命令行参数
	cfToken      string
	serverName   string
	isLocal      bool
	installTrust *bool
	fromHost     string
	toURL        string
	domain       string
	subdomain    string
	ports        string
	host         string
	routeID      string
)

// rootCmd 根命令 - FastCaddy CLI 工具的主入口
var rootCmd = &cobra.Command{
	Use:   "fastcaddy",
	Short: "FastCaddy - Caddy 配置管理工具",
	Long: `FastCaddy 是一个用于简化 Caddy Web 服务器配置的命令行工具。
它提供了简单的命令来设置反向代理、SSL 证书和路由规则。

示例:
  fastcaddy setup --local                          # 设置本地开发环境
  fastcaddy setup --cf-token $CADDY_CF_TOKEN       # 设置生产环境
  fastcaddy add-proxy --from api.example.com --to localhost:8080
  fastcaddy add-wildcard --domain example.com`,
}

// setupCmd 设置命令 - 初始化 Caddy 基本配置
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "设置 Caddy 基本配置",
	Long: `初始化 Caddy 的基本配置，包括 SSL/TLS 设置和 HTTP 服务器配置。

可以配置为本地开发环境（使用内部证书）或生产环境（使用 ACME/Let's Encrypt）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fc := fastcaddy.New()

		// 如果没有提供 CF Token，尝试从环境变量获取
		if cfToken == "" && !isLocal {
			cfToken = utils.GetCloudflareToken()
		}

		// 执行 Caddy 设置
		fmt.Printf("正在设置 Caddy 配置...\n")
		if isLocal {
			fmt.Printf("配置类型: 本地开发环境（内部证书）\n")
		} else {
			fmt.Printf("配置类型: 生产环境（ACME 证书）\n")
			if cfToken != "" {
				fmt.Printf("使用 Cloudflare DNS 挑战\n")
			}
		}

		err := fc.SetupCaddy(cfToken, serverName, isLocal, installTrust)
		if err != nil {
			return fmt.Errorf("设置 Caddy 失败: %w", err)
		}

		fmt.Printf("✓ Caddy 配置设置成功\n")
		return nil
	},
}

// addProxyCmd 添加反向代理命令
var addProxyCmd = &cobra.Command{
	Use:   "add-proxy",
	Short: "添加反向代理",
	Long: `为指定主机添加反向代理配置。

示例:
  fastcaddy add-proxy --from api.example.com --to localhost:8080
  fastcaddy add-proxy --from web.example.com --to 127.0.0.1:3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if fromHost == "" || toURL == "" {
			return fmt.Errorf("必须指定 --from 和 --to 参数")
		}

		if !utils.ValidateHost(fromHost) {
			return fmt.Errorf("无效的主机名: %s", fromHost)
		}

		if !utils.ValidateURL(toURL) {
			return fmt.Errorf("无效的目标 URL: %s", toURL)
		}

		fc := fastcaddy.New()

		fmt.Printf("正在添加反向代理: %s -> %s\n", fromHost, toURL)
		err := fc.AddReverseProxy(fromHost, toURL)
		if err != nil {
			return fmt.Errorf("添加反向代理失败: %w", err)
		}

		fmt.Printf("✓ 反向代理添加成功\n")
		return nil
	},
}

// delProxyCmd 删除反向代理命令
var delProxyCmd = &cobra.Command{
	Use:   "del-proxy",
	Short: "删除反向代理",
	Long: `删除指定 ID 的反向代理配置。

示例:
  fastcaddy del-proxy --id api.example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if routeID == "" {
			return fmt.Errorf("必须指定 --id 参数")
		}

		fc := fastcaddy.New()

		if !fc.HasID(routeID) {
			return fmt.Errorf("路由 ID '%s' 不存在", routeID)
		}

		fmt.Printf("正在删除路由: %s\n", routeID)
		err := fc.DeleteRoute(routeID)
		if err != nil {
			return fmt.Errorf("删除路由失败: %w", err)
		}

		fmt.Printf("✓ 路由删除成功\n")
		return nil
	},
}

// addWildcardCmd 添加通配符命令
var addWildcardCmd = &cobra.Command{
	Use:   "add-wildcard",
	Short: "添加通配符子域名",
	Long: `为指定域名添加通配符子域名支持。

示例:
  fastcaddy add-wildcard --domain example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if domain == "" {
			return fmt.Errorf("必须指定 --domain 参数")
		}

		fc := fastcaddy.New()

		fmt.Printf("正在添加通配符路由: *.%s\n", domain)
		err := fc.AddWildcardRoute(domain)
		if err != nil {
			return fmt.Errorf("添加通配符路由失败: %w", err)
		}

		fmt.Printf("✓ 通配符路由添加成功\n")
		return nil
	},
}

// addSubProxyCmd 添加子域名反向代理命令
var addSubProxyCmd = &cobra.Command{
	Use:   "add-sub-proxy",
	Short: "添加子域名反向代理",
	Long: `为通配符域名下的特定子域名添加反向代理。

示例:
  fastcaddy add-sub-proxy --domain example.com --subdomain api --ports 8080 --host localhost
  fastcaddy add-sub-proxy --domain example.com --subdomain web --ports 3000,3001`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if domain == "" || subdomain == "" || ports == "" {
			return fmt.Errorf("必须指定 --domain, --subdomain 和 --ports 参数")
		}

		if host == "" {
			host = "localhost"
		}

		// 解析端口列表
		portList := strings.Split(ports, ",")
		for i, port := range portList {
			portList[i] = strings.TrimSpace(port)
			// 验证端口是否为数字
			if _, err := strconv.Atoi(portList[i]); err != nil {
				return fmt.Errorf("无效的端口号: %s", portList[i])
			}
		}

		fc := fastcaddy.New()

		fmt.Printf("正在添加子域名反向代理: %s.%s -> %s:%s\n", subdomain, domain, host, ports)
		err := fc.AddSubReverseProxy(domain, subdomain, portList, host)
		if err != nil {
			return fmt.Errorf("添加子域名反向代理失败: %w", err)
		}

		fmt.Printf("✓ 子域名反向代理添加成功\n")
		return nil
	},
}

// statusCmd 状态命令
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看 Caddy 配置状态",
	Long:  `显示当前 Caddy 配置的状态信息。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fc := fastcaddy.New()

		// 检查各个配置路径是否存在
		fmt.Printf("Caddy 配置状态:\n")
		fmt.Printf("===============\n")

		if fc.HasPath("/apps/http/servers") {
			fmt.Printf("✓ HTTP 服务器: 已配置\n")
		} else {
			fmt.Printf("✗ HTTP 服务器: 未配置\n")
		}

		if fc.HasPath("/apps/tls/automation") {
			fmt.Printf("✓ TLS 自动化: 已配置\n")
		} else {
			fmt.Printf("✗ TLS 自动化: 未配置\n")
		}

		if fc.HasPath("/apps/pki") {
			fmt.Printf("✓ PKI 配置: 已配置\n")
		} else {
			fmt.Printf("✗ PKI 配置: 未配置\n")
		}

		return nil
	},
}

func init() {
	// 设置命令参数
	setupCmd.Flags().StringVar(&cfToken, "cf-token", "", "Cloudflare API 令牌（用于 ACME DNS 挑战）")
	setupCmd.Flags().StringVar(&serverName, "server", "srv0", "服务器名称")
	setupCmd.Flags().BoolVar(&isLocal, "local", false, "是否为本地开发环境（使用内部证书）")
	
	// installTrust 参数需要特殊处理，因为它是一个 *bool
	var installTrustFlag bool
	setupCmd.Flags().BoolVar(&installTrustFlag, "install-trust", false, "是否安装根证书到系统信任存储")
	setupCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("install-trust") {
			installTrust = &installTrustFlag
		}
		return nil
	}

	addProxyCmd.Flags().StringVar(&fromHost, "from", "", "源主机名（必需）")
	addProxyCmd.Flags().StringVar(&toURL, "to", "", "目标 URL（必需）")
	addProxyCmd.MarkFlagRequired("from")
	addProxyCmd.MarkFlagRequired("to")

	delProxyCmd.Flags().StringVar(&routeID, "id", "", "路由 ID（必需）")
	delProxyCmd.MarkFlagRequired("id")

	addWildcardCmd.Flags().StringVar(&domain, "domain", "", "域名（必需）")
	addWildcardCmd.MarkFlagRequired("domain")

	addSubProxyCmd.Flags().StringVar(&domain, "domain", "", "域名（必需）")
	addSubProxyCmd.Flags().StringVar(&subdomain, "subdomain", "", "子域名（必需）")
	addSubProxyCmd.Flags().StringVar(&ports, "ports", "", "端口列表，用逗号分隔（必需）")
	addSubProxyCmd.Flags().StringVar(&host, "host", "localhost", "目标主机")
	addSubProxyCmd.MarkFlagRequired("domain")
	addSubProxyCmd.MarkFlagRequired("subdomain")
	addSubProxyCmd.MarkFlagRequired("ports")

	// 添加子命令到根命令
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(addProxyCmd)
	rootCmd.AddCommand(delProxyCmd)
	rootCmd.AddCommand(addWildcardCmd)
	rootCmd.AddCommand(addSubProxyCmd)
	rootCmd.AddCommand(statusCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}