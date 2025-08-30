#!/usr/bin/env python3
"""
FastCaddy Python 版本高级示例
展示域名管理、配置检查、删除操作等编程接口的使用

使用说明：
1. 确保 Caddy 服务器运行在 localhost:2019
2. 安装 Python 依赖: pip install fastcore httpx
3. 运行此脚本: python3 advanced_example.py
"""

import sys
import os
import time

# 添加 fastcaddy 模块路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from fastcaddy.core import *

def print_separator(title):
    """打印分隔线和标题"""
    print(f"\n{'='*50}")
    print(f" {title}")
    print(f"{'='*50}")

def check_domain_status(domain):
    """检查域名配置状态"""
    if has_id(domain):
        print(f"✅ 域名 {domain} 已配置")
        return True
    else:
        print(f"❌ 域名 {domain} 未配置")
        return False

def safe_add_domain(domain, target):
    """安全添加域名（检查现有配置）"""
    print(f"\n🔍 检查域名 {domain}...")
    
    if has_id(domain):
        print(f"⚠️  域名 {domain} 已存在，将先删除")
        try:
            del_id(domain)
            print(f"🗑️  成功删除现有域名 {domain}")
            time.sleep(0.1)  # 短暂延迟
        except Exception as e:
            print(f"❌ 删除现有域名失败: {e}")
            return False
    
    print(f"➕ 添加域名 {domain} -> {target}")
    try:
        add_reverse_proxy(domain, target)
        print(f"✅ 成功添加域名 {domain}")
        return True
    except Exception as e:
        print(f"❌ 添加域名失败: {e}")
        return False

def safe_delete_domain(domain):
    """安全删除域名（检查后删除）"""
    print(f"\n🔍 检查域名 {domain}...")
    
    if not has_id(domain):
        print(f"⚠️  域名 {domain} 不存在，无需删除")
        return True
    
    print(f"🗑️  删除域名 {domain}")
    try:
        del_id(domain)
        # 验证删除结果
        if has_id(domain):
            print(f"❌ 域名 {domain} 删除失败，仍然存在")
            return False
        print(f"✅ 成功删除域名 {domain}")
        return True
    except Exception as e:
        print(f"❌ 删除域名失败: {e}")
        return False

def batch_manage_domains(domains, action):
    """批量管理域名"""
    print(f"\n📋 批量{action}域名...")
    results = {}
    
    for domain, target in domains.items():
        if action == "添加":
            results[domain] = safe_add_domain(domain, target)
        elif action == "删除":
            results[domain] = safe_delete_domain(domain)
        elif action == "检查":
            results[domain] = check_domain_status(domain)
    
    return results

def main():
    print("FastCaddy Python 版本高级示例程序")
    print("===================================")

    # 初始化环境
    print_separator("🚀 初始化环境")
    try:
        setup_caddy(local=True, install_trust=True)
        print("✅ 环境初始化完成")
    except Exception as e:
        print(f"⚠️  初始化失败: {e} (可能是因为 Caddy 服务器未运行)")

    # 示例域名配置
    domains = {
        "api.test.com": "localhost:8080",
        "web.test.com": "localhost:3000", 
        "admin.test.com": "localhost:9000",
        "cache.test.com": "localhost:6379",
        "db.test.com": "localhost:5432",
    }

    # 场景1: 检查所有域名状态
    print_separator("📊 场景1: 检查当前域名状态")
    batch_manage_domains(domains, "检查")

    # 场景2: 批量添加域名
    print_separator("📝 场景2: 批量添加域名")
    results = batch_manage_domains(domains, "添加")
    success_count = sum(1 for success in results.values() if success)
    print(f"\n📊 添加结果: {success_count}/{len(domains)} 个域名添加成功")

    # 场景3: 验证添加结果
    print_separator("✅ 场景3: 验证添加结果")
    batch_manage_domains(domains, "检查")

    # 场景4: 更新域名配置
    print_separator("🔄 场景4: 更新域名配置")
    updated_domains = {
        "api.test.com": "localhost:8081",  # 端口变更
        "web.test.com": "localhost:3001",  # 端口变更
    }
    batch_manage_domains(updated_domains, "添加")

    # 场景5: 检查配置详情
    print_separator("📊 场景5: 检查配置详情")
    
    # 检查HTTP服务器配置
    if has_path("/apps/http"):
        print("✅ HTTP 应用已配置")
        try:
            http_config = gcfg("/apps/http/servers")
            print(f"   HTTP服务器数量: {len(http_config)}")
        except Exception as e:
            print(f"   ❌ 获取HTTP配置失败: {e}")
    else:
        print("❌ HTTP 应用未配置")

    # 检查TLS配置
    if has_path("/apps/tls"):
        print("✅ TLS 应用已配置")
    else:
        print("⚠️  TLS 应用未配置")

    # 场景6: 通配符域名管理
    print_separator("🌟 场景6: 通配符域名管理")
    wildcard_domain = "dev.local"
    wildcard_id = f"wildcard-{wildcard_domain}"
    
    # 添加通配符路由
    if not has_id(wildcard_id):
        try:
            add_wildcard_route(wildcard_domain)
            print(f"✅ 添加通配符域名 *.{wildcard_domain}")
        except Exception as e:
            print(f"❌ 添加通配符域名失败: {e}")
    else:
        print(f"⚠️  通配符域名 *.{wildcard_domain} 已存在")

    # 添加子域名配置
    sub_domains = {
        "app": "8090",
        "api": "8091", 
        "admin": "8092",
    }

    print("\n📋 添加子域名配置...")
    for subdomain, port in sub_domains.items():
        full_domain = f"{subdomain}.{wildcard_domain}"
        if not has_id(full_domain):
            try:
                add_sub_reverse_proxy(wildcard_domain, subdomain, port)
                print(f"✅ 添加子域名 {full_domain} -> localhost:{port}")
            except Exception as e:
                print(f"❌ 添加子域名 {full_domain} 失败: {e}")
        else:
            print(f"⚠️  子域名 {full_domain} 已存在")

    # 检查所有子域名状态
    print("\n🔍 检查子域名状态...")
    for subdomain in sub_domains:
        full_domain = f"{subdomain}.{wildcard_domain}"
        check_domain_status(full_domain)

    # 场景7: 选择性删除域名
    print_separator("🗑️ 场景7: 删除部分域名")
    domains_to_delete = {
        "cache.test.com": "",
        "db.test.com": "",
    }
    delete_results = batch_manage_domains(domains_to_delete, "删除")
    success_count = sum(1 for success in delete_results.values() if success)
    print(f"\n📊 删除结果: {success_count}/{len(domains_to_delete)} 个域名删除成功")

    # 场景8: 最终状态检查
    print_separator("📈 场景8: 最终状态检查")
    final_status = batch_manage_domains(domains, "检查")
    active_domains = [domain for domain, exists in final_status.items() if exists]
    print(f"\n📊 活跃域名: {len(active_domains)}/{len(domains)}")
    for domain in active_domains:
        print(f"   🟢 {domain} (活跃)")

    # 场景9: 清理所有配置
    print_separator("🧹 场景9: 清理所有配置")
    
    # 删除子域名
    print("删除子域名...")
    for subdomain in sub_domains:
        full_domain = f"{subdomain}.{wildcard_domain}"
        safe_delete_domain(full_domain)
    
    # 删除通配符域名
    print("删除通配符域名...")
    safe_delete_domain(wildcard_id)
    
    # 删除剩余普通域名
    print("删除剩余普通域名...")
    remaining_domains = {
        "api.test.com": "",
        "web.test.com": "",
        "admin.test.com": "",
    }
    batch_manage_domains(remaining_domains, "删除")

    print_separator("🎉 高级示例程序完成!")
    print("\n💡 Python 编程接口使用技巧:")
    print("===============================")
    print("1. 使用 has_id(domain) 检查域名是否已配置")
    print("2. 使用 has_path(path) 检查配置路径是否存在") 
    print("3. 使用 del_id(id) 删除指定的路由配置")
    print("4. 使用 gcfg(path) 获取详细配置信息")
    print("5. 使用 add_reverse_proxy(from, to) 添加反向代理")
    print("6. 使用 add_wildcard_route(domain) 添加通配符路由")
    print("7. 使用 add_sub_reverse_proxy(domain, sub, port) 添加子域名")
    print("8. 批量操作前建议先检查状态，避免重复操作")
    print("9. 删除操作后建议验证结果，确保操作成功")

if __name__ == "__main__":
    main()