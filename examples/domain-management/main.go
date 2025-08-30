package main

import (
	"fmt"
	"log"
	"time"

	"github.com/youfun/fastcaddy"
)

// 示例：域名管理专用示例 - 展示如何检查、添加、删除域名配置
func main() {
	fmt.Println("FastCaddy 域名管理示例")
	fmt.Println("=======================")

	// 创建客户端
	fc := fastcaddy.New()

	// 初始化
	fmt.Println("\n🚀 初始化环境...")
	installTrust := true
	err := fc.SetupCaddy("", "srv0", true, &installTrust)
	if err != nil {
		log.Printf("初始化失败: %v", err)
	}

	// 域名管理函数
	type DomainManager struct {
		client *fastcaddy.FastCaddy
	}

	dm := &DomainManager{client: fc}

	// 检查域名是否已配置
	checkDomain := func(domain string) bool {
		if dm.client.HasID(domain) {
			fmt.Printf("✅ 域名 %s 已配置\n", domain)
			return true
		} else {
			fmt.Printf("❌ 域名 %s 未配置\n", domain)
			return false
		}
	}

	// 安全添加域名（检查后添加）
	safeAddDomain := func(domain, target string) error {
		fmt.Printf("\n🔍 检查域名 %s...\n", domain)
		if dm.client.HasID(domain) {
			fmt.Printf("⚠️  域名 %s 已存在，将先删除\n", domain)
			err := dm.client.DeleteRoute(domain)
			if err != nil {
				return fmt.Errorf("删除现有域名失败: %v", err)
			}
			fmt.Printf("🗑️  成功删除现有域名 %s\n", domain)
			time.Sleep(100 * time.Millisecond) // 短暂延迟
		}
		
		fmt.Printf("➕ 添加域名 %s -> %s\n", domain, target)
		err := dm.client.AddReverseProxy(domain, target)
		if err != nil {
			return fmt.Errorf("添加域名失败: %v", err)
		}
		fmt.Printf("✅ 成功添加域名 %s\n", domain)
		return nil
	}

	// 安全删除域名（检查后删除）
	safeDeleteDomain := func(domain string) error {
		fmt.Printf("\n🔍 检查域名 %s...\n", domain)
		if !dm.client.HasID(domain) {
			fmt.Printf("⚠️  域名 %s 不存在，无需删除\n", domain)
			return nil
		}
		
		fmt.Printf("🗑️  删除域名 %s\n", domain)
		err := dm.client.DeleteRoute(domain)
		if err != nil {
			return fmt.Errorf("删除域名失败: %v", err)
		}
		
		// 验证删除结果
		if dm.client.HasID(domain) {
			return fmt.Errorf("域名 %s 删除失败，仍然存在", domain)
		}
		fmt.Printf("✅ 成功删除域名 %s\n", domain)
		return nil
	}

	// 批量管理域名
	manageDomains := func(domains map[string]string, action string) {
		fmt.Printf("\n📋 批量%s域名...\n", action)
		for domain, target := range domains {
			switch action {
			case "添加":
				if err := safeAddDomain(domain, target); err != nil {
					log.Printf("❌ %s: %v", domain, err)
				}
			case "删除":
				if err := safeDeleteDomain(domain); err != nil {
					log.Printf("❌ %s: %v", domain, err)
				}
			case "检查":
				checkDomain(domain)
			}
		}
	}

	// 示例域名配置
	domains := map[string]string{
		"api.test.com":    "localhost:8080",
		"web.test.com":    "localhost:3000", 
		"admin.test.com":  "localhost:9000",
		"cache.test.com":  "localhost:6379",
		"db.test.com":     "localhost:5432",
	}

	// 场景1: 检查所有域名状态
	fmt.Println("\n📊 场景1: 检查当前域名状态")
	manageDomains(domains, "检查")

	// 场景2: 批量添加域名
	fmt.Println("\n📝 场景2: 批量添加域名")
	manageDomains(domains, "添加")

	// 场景3: 验证添加结果
	fmt.Println("\n✅ 场景3: 验证添加结果")
	manageDomains(domains, "检查")

	// 场景4: 更新域名配置（重新添加）
	fmt.Println("\n🔄 场景4: 更新域名配置")
	updatedDomains := map[string]string{
		"api.test.com": "localhost:8081",  // 端口变更
		"web.test.com": "localhost:3001",  // 端口变更
	}
	manageDomains(updatedDomains, "添加")

	// 场景5: 选择性删除域名
	fmt.Println("\n🗑️  场景5: 删除部分域名")
	domainsToDelete := map[string]string{
		"cache.test.com": "",
		"db.test.com":    "",
	}
	manageDomains(domainsToDelete, "删除")

	// 场景6: 最终状态检查
	fmt.Println("\n📈 场景6: 最终状态检查")
	manageDomains(domains, "检查")

	// 高级功能：通配符域名管理
	fmt.Println("\n🌟 高级功能: 通配符域名管理")
	wildcardDomain := "dev.local"
	
	// 添加通配符路由
	if !checkDomain(fmt.Sprintf("wildcard-%s", wildcardDomain)) {
		err := fc.AddWildcardRoute(wildcardDomain)
		if err != nil {
			log.Printf("添加通配符域名失败: %v", err)
		} else {
			fmt.Printf("✅ 添加通配符域名 *.%s\n", wildcardDomain)
		}
	}

	// 添加子域名配置
	subDomains := map[string]string{
		"app":      "8090",
		"api":      "8091", 
		"admin":    "8092",
	}

	fmt.Println("\n📋 添加子域名配置...")
	for subdomain, port := range subDomains {
		fullDomain := fmt.Sprintf("%s.%s", subdomain, wildcardDomain)
		if !checkDomain(fullDomain) {
			err := fc.AddSubReverseProxy(wildcardDomain, subdomain, []string{port}, "localhost")
			if err != nil {
				log.Printf("添加子域名 %s 失败: %v", fullDomain, err)
			} else {
				fmt.Printf("✅ 添加子域名 %s -> localhost:%s\n", fullDomain, port)
			}
		}
	}

	// 检查所有子域名状态
	fmt.Println("\n🔍 检查子域名状态...")
	for subdomain := range subDomains {
		fullDomain := fmt.Sprintf("%s.%s", subdomain, wildcardDomain)
		checkDomain(fullDomain)
	}

	// 清理所有配置
	fmt.Println("\n🧹 清理所有配置...")
	
	// 删除子域名
	for subdomain := range subDomains {
		fullDomain := fmt.Sprintf("%s.%s", subdomain, wildcardDomain)
		if err := safeDeleteDomain(fullDomain); err != nil {
			log.Printf("删除子域名失败: %v", err)
		}
	}
	
	// 删除通配符域名
	wildcardID := fmt.Sprintf("wildcard-%s", wildcardDomain)
	if err := safeDeleteDomain(wildcardID); err != nil {
		log.Printf("删除通配符域名失败: %v", err)
	}
	
	// 删除剩余普通域名
	remainingDomains := map[string]string{
		"api.test.com":   "",
		"web.test.com":   "",
		"admin.test.com": "",
	}
	manageDomains(remainingDomains, "删除")

	fmt.Println("\n🎉 域名管理示例完成!")
	fmt.Println("\n💡 域名管理最佳实践:")
	fmt.Println("======================")
	fmt.Println("1. 添加前检查域名是否已存在 (HasID)")
	fmt.Println("2. 更新域名时先删除再添加")
	fmt.Println("3. 删除后验证操作是否成功")
	fmt.Println("4. 批量操作时处理每个操作的错误")
	fmt.Println("5. 使用通配符域名可以简化子域名管理")
	fmt.Println("6. 定期检查配置状态确保一致性")
}