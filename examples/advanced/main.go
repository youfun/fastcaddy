package main

import (
	"fmt"
	"log"

	"github.com/youfun/fastcaddy"
)

// 示例：展示 FastCaddy Go 版本的高级编程接口使用
// 包括域名管理、检查配置状态、删除操作等
func main() {
	fmt.Println("FastCaddy Go 版本高级示例程序")
	fmt.Println("================================")

	// 创建 FastCaddy 客户端实例
	fc := fastcaddy.New()

	// 初始化基础配置
	fmt.Println("\n🚀 初始化环境...")
	installTrust := true
	err := fc.SetupCaddy("", "srv0", true, &installTrust)
	if err != nil {
		log.Printf("初始化失败: %v (可能是因为 Caddy 服务器未运行)", err)
	} else {
		fmt.Println("✅ 环境初始化完成")
	}

	// 示例域名列表
	domains := []string{"api.example.com", "web.example.com", "admin.example.com"}
	ports := []string{"8080", "3000", "9000"}

	// 1. 批量添加反向代理
	fmt.Println("\n📝 1. 批量添加反向代理...")
	for i, domain := range domains {
		target := fmt.Sprintf("localhost:%s", ports[i])
		err := fc.AddReverseProxy(domain, target)
		if err != nil {
			log.Printf("   ❌ 添加 %s -> %s 失败: %v", domain, target, err)
		} else {
			fmt.Printf("   ✅ 添加 %s -> %s 成功\n", domain, target)
		}
	}

	// 2. 检查域名配置状态
	fmt.Println("\n🔍 2. 检查域名配置状态...")
	for _, domain := range domains {
		if fc.HasID(domain) {
			fmt.Printf("   ✅ %s 已配置\n", domain)
		} else {
			fmt.Printf("   ❌ %s 未配置\n", domain)
		}
	}

	// 3. 添加通配符域名和子域名
	fmt.Println("\n🌟 3. 添加通配符域名和子域名...")
	wildcardDomain := "dev.local"
	
	// 添加通配符路由
	err = fc.AddWildcardRoute(wildcardDomain)
	if err != nil {
		log.Printf("   ❌ 添加通配符域名 *.%s 失败: %v", wildcardDomain, err)
	} else {
		fmt.Printf("   ✅ 添加通配符域名 *.%s 成功\n", wildcardDomain)
	}

	// 添加子域名
	subdomains := []string{"app", "dashboard", "monitoring"}
	subPorts := []string{"8081", "8082", "8083"}
	
	for i, subdomain := range subdomains {
		err := fc.AddSubReverseProxy(wildcardDomain, subdomain, []string{subPorts[i]}, "localhost")
		if err != nil {
			log.Printf("   ❌ 添加子域名 %s.%s 失败: %v", subdomain, wildcardDomain, err)
		} else {
			fmt.Printf("   ✅ 添加子域名 %s.%s -> localhost:%s 成功\n", subdomain, wildcardDomain, subPorts[i])
		}
	}

	// 4. 检查子域名配置状态
	fmt.Println("\n🔍 4. 检查子域名配置状态...")
	for _, subdomain := range subdomains {
		fullDomain := fmt.Sprintf("%s.%s", subdomain, wildcardDomain)
		if fc.HasID(fullDomain) {
			fmt.Printf("   ✅ %s 已配置\n", fullDomain)
		} else {
			fmt.Printf("   ❌ %s 未配置\n", fullDomain)
		}
	}

	// 5. 获取详细配置信息
	fmt.Println("\n📊 5. 获取配置信息...")
	
	// 检查HTTP服务器配置
	if fc.HasPath("/apps/http") {
		config, err := fc.GetConfig("/apps/http/servers")
		if err != nil {
			log.Printf("   ❌ 获取HTTP服务器配置失败: %v", err)
		} else {
			fmt.Printf("   ✅ HTTP服务器配置: %d 个服务器\n", len(config))
		}
	}

	// 检查TLS配置
	if fc.HasPath("/apps/tls") {
		fmt.Println("   ✅ TLS 配置已启用")
	} else {
		fmt.Println("   ⚠️  TLS 配置未启用")
	}

	// 6. 删除指定域名配置
	fmt.Println("\n🗑️ 6. 删除域名配置...")
	
	// 删除第一个域名
	domainToDelete := domains[0]
	if fc.HasID(domainToDelete) {
		err := fc.DeleteRoute(domainToDelete)
		if err != nil {
			log.Printf("   ❌ 删除 %s 失败: %v", domainToDelete, err)
		} else {
			fmt.Printf("   ✅ 删除 %s 成功\n", domainToDelete)
		}
	} else {
		fmt.Printf("   ⚠️  %s 不存在，无需删除\n", domainToDelete)
	}

	// 7. 验证删除结果
	fmt.Println("\n✅ 7. 验证删除结果...")
	if fc.HasID(domainToDelete) {
		fmt.Printf("   ❌ %s 仍然存在（删除失败）\n", domainToDelete)
	} else {
		fmt.Printf("   ✅ %s 已成功删除\n", domainToDelete)
	}

	// 8. 批量删除子域名
	fmt.Println("\n🧹 8. 批量删除子域名...")
	for _, subdomain := range subdomains {
		fullDomain := fmt.Sprintf("%s.%s", subdomain, wildcardDomain)
		if fc.HasID(fullDomain) {
			err := fc.DeleteRoute(fullDomain)
			if err != nil {
				log.Printf("   ❌ 删除 %s 失败: %v", fullDomain, err)
			} else {
				fmt.Printf("   ✅ 删除 %s 成功\n", fullDomain)
			}
		}
	}

	// 9. 最终状态检查
	fmt.Println("\n📈 9. 最终状态检查...")
	allDomains := append(domains, fmt.Sprintf("wildcard-%s", wildcardDomain))
	for _, subdomain := range subdomains {
		allDomains = append(allDomains, fmt.Sprintf("%s.%s", subdomain, wildcardDomain))
	}
	
	activeCount := 0
	for _, domain := range allDomains {
		if fc.HasID(domain) {
			activeCount++
			fmt.Printf("   🟢 %s (活跃)\n", domain)
		} else {
			fmt.Printf("   🔴 %s (已删除)\n", domain)
		}
	}
	
	fmt.Printf("\n📊 总结: %d/%d 个域名配置仍然活跃\n", activeCount, len(allDomains))

	fmt.Println("\n🎉 高级示例程序完成!")
	fmt.Println("\n💡 编程接口使用技巧:")
	fmt.Println("====================")
	fmt.Println("1. 使用 fc.HasID(domain) 检查域名是否已配置")
	fmt.Println("2. 使用 fc.HasPath(path) 检查配置路径是否存在") 
	fmt.Println("3. 使用 fc.DeleteRoute(id) 删除指定的路由配置")
	fmt.Println("4. 使用 fc.GetConfig(path) 获取详细配置信息")
	fmt.Println("5. 批量操作前建议先检查状态，避免重复操作")
	fmt.Println("6. 删除操作后建议验证结果，确保操作成功")
}