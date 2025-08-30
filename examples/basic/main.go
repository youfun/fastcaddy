package main

import (
	"fmt"
	"log"

	"github.com/youfun/fastcaddy"
)

// 示例：展示 FastCaddy Go 版本的基本使用
func main() {
	fmt.Println("FastCaddy Go 版本示例程序")
	fmt.Println("==========================")

	// 创建 FastCaddy 客户端实例
	fc := fastcaddy.New()

	// 示例 1: 设置本地开发环境
	fmt.Println("\n1. 设置本地开发环境...")
	installTrust := true
	err := fc.SetupCaddy("", "srv0", true, &installTrust)
	if err != nil {
		log.Printf("设置 Caddy 失败: %v", err)
		log.Println("注意: 这可能是因为没有运行 Caddy 服务器")
	} else {
		fmt.Println("✅ Caddy 配置设置成功!")
	}

	// 示例 2: 添加反向代理
	fmt.Println("\n2. 添加反向代理...")
	err = fc.AddReverseProxy("api.localhost", "localhost:8080")
	if err != nil {
		log.Printf("添加反向代理失败: %v", err)
	} else {
		fmt.Println("✅ 反向代理添加成功: api.localhost -> localhost:8080")
	}

	// 示例 3: 添加通配符域名
	fmt.Println("\n3. 添加通配符域名...")
	err = fc.AddWildcardRoute("localhost")
	if err != nil {
		log.Printf("添加通配符域名失败: %v", err)
	} else {
		fmt.Println("✅ 通配符域名添加成功: *.localhost")
	}

	// 示例 4: 添加子域名反向代理
	fmt.Println("\n4. 添加子域名反向代理...")
	err = fc.AddSubReverseProxy("localhost", "web", []string{"3000"}, "localhost")
	if err != nil {
		log.Printf("添加子域名反向代理失败: %v", err)
	} else {
		fmt.Println("✅ 子域名反向代理添加成功: web.localhost -> localhost:3000")
	}

	// 示例 5: 检查配置状态
	fmt.Println("\n5. 检查配置状态...")
	if fc.HasPath("/apps/http/servers") {
		fmt.Println("✅ HTTP 服务器已配置")
	} else {
		fmt.Println("❌ HTTP 服务器未配置")
	}

	if fc.HasPath("/apps/tls/automation") {
		fmt.Println("✅ TLS 自动化已配置")
	} else {
		fmt.Println("❌ TLS 自动化未配置")
	}

	// 示例 6: 获取配置信息（如果可用）
	fmt.Println("\n6. 获取配置信息...")
	config, err := fc.GetConfig("/")
	if err != nil {
		log.Printf("获取配置失败: %v", err)
	} else {
		fmt.Printf("根配置包含 %d 个应用\n", len(config))
		for appName := range config {
			fmt.Printf("- 应用: %s\n", appName)
		}
	}

	fmt.Println("\n示例程序完成!")
	fmt.Println("\n使用说明:")
	fmt.Println("=========")
	fmt.Println("1. 确保 Caddy 服务器正在运行 (通常在 localhost:2019)")
	fmt.Println("2. 使用 './fastcaddy-go status' 检查配置状态")
	fmt.Println("3. 使用 './fastcaddy-go setup --local' 初始化本地环境")
	fmt.Println("4. 使用 './fastcaddy-go add-proxy --from example.localhost --to localhost:8080' 添加代理")
}