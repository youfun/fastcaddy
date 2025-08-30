#!/usr/bin/env python3
"""
FastCaddy Python 版本域名管理专用示例
展示如何检查、添加、删除域名配置的最佳实践

使用说明：
1. 确保 Caddy 服务器运行在 localhost:2019
2. 安装 Python 依赖: pip install fastcore httpx
3. 运行此脚本: python3 domain_management.py
"""

import sys
import os
import time
from typing import Dict, List, Optional

# 添加 fastcaddy 模块路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from fastcaddy.core import *

class DomainManager:
    """域名管理器 - 提供安全的域名操作方法"""
    
    def __init__(self):
        self.setup_environment()
    
    def setup_environment(self):
        """初始化环境"""
        try:
            setup_caddy(local=True, install_trust=True)
            print("✅ FastCaddy 环境初始化完成")
        except Exception as e:
            print(f"⚠️  环境初始化失败: {e}")
    
    def check_domain(self, domain: str) -> bool:
        """检查域名是否已配置"""
        exists = has_id(domain)
        status = "已配置" if exists else "未配置"
        print(f"{'✅' if exists else '❌'} 域名 {domain} {status}")
        return exists
    
    def add_domain(self, domain: str, target: str, force: bool = False) -> bool:
        """添加域名配置"""
        print(f"\n🔍 检查域名 {domain}...")
        
        if has_id(domain):
            if force:
                print(f"⚠️  域名 {domain} 已存在，强制更新模式")
                if not self.delete_domain(domain, silent=True):
                    return False
                time.sleep(0.1)  # 短暂延迟确保删除完成
            else:
                print(f"❌ 域名 {domain} 已存在，使用 force=True 强制更新")
                return False
        
        try:
            print(f"➕ 添加域名 {domain} -> {target}")
            add_reverse_proxy(domain, target)
            print(f"✅ 成功添加域名 {domain}")
            return True
        except Exception as e:
            print(f"❌ 添加域名失败: {e}")
            return False
    
    def delete_domain(self, domain: str, silent: bool = False) -> bool:
        """删除域名配置"""
        if not silent:
            print(f"\n🔍 检查域名 {domain}...")
        
        if not has_id(domain):
            if not silent:
                print(f"⚠️  域名 {domain} 不存在，无需删除")
            return True
        
        try:
            if not silent:
                print(f"🗑️  删除域名 {domain}")
            del_id(domain)
            
            # 验证删除结果
            if has_id(domain):
                if not silent:
                    print(f"❌ 域名 {domain} 删除失败，仍然存在")
                return False
            
            if not silent:
                print(f"✅ 成功删除域名 {domain}")
            return True
        except Exception as e:
            if not silent:
                print(f"❌ 删除域名失败: {e}")
            return False
    
    def update_domain(self, domain: str, new_target: str) -> bool:
        """更新域名配置"""
        print(f"\n🔄 更新域名 {domain} -> {new_target}")
        return self.add_domain(domain, new_target, force=True)
    
    def batch_check(self, domains: List[str]) -> Dict[str, bool]:
        """批量检查域名状态"""
        print(f"\n📊 批量检查 {len(domains)} 个域名状态...")
        results = {}
        for domain in domains:
            results[domain] = self.check_domain(domain)
        return results
    
    def batch_add(self, domain_configs: Dict[str, str], force: bool = False) -> Dict[str, bool]:
        """批量添加域名"""
        print(f"\n📝 批量添加 {len(domain_configs)} 个域名...")
        results = {}
        for domain, target in domain_configs.items():
            results[domain] = self.add_domain(domain, target, force)
        return results
    
    def batch_delete(self, domains: List[str]) -> Dict[str, bool]:
        """批量删除域名"""
        print(f"\n🗑️  批量删除 {len(domains)} 个域名...")
        results = {}
        for domain in domains:
            results[domain] = self.delete_domain(domain)
        return results
    
    def add_wildcard_domain(self, base_domain: str) -> bool:
        """添加通配符域名"""
        wildcard_id = f"wildcard-{base_domain}"
        
        if has_id(wildcard_id):
            print(f"⚠️  通配符域名 *.{base_domain} 已存在")
            return False
        
        try:
            add_wildcard_route(base_domain)
            print(f"✅ 添加通配符域名 *.{base_domain}")
            return True
        except Exception as e:
            print(f"❌ 添加通配符域名失败: {e}")
            return False
    
    def add_subdomain(self, base_domain: str, subdomain: str, port: str) -> bool:
        """添加子域名"""
        full_domain = f"{subdomain}.{base_domain}"
        
        if has_id(full_domain):
            print(f"⚠️  子域名 {full_domain} 已存在")
            return False
        
        try:
            add_sub_reverse_proxy(base_domain, subdomain, port)
            print(f"✅ 添加子域名 {full_domain} -> localhost:{port}")
            return True
        except Exception as e:
            print(f"❌ 添加子域名失败: {e}")
            return False
    
    def get_system_status(self) -> Dict[str, bool]:
        """获取系统状态"""
        status = {
            "http_app": has_path("/apps/http"),
            "tls_app": has_path("/apps/tls"),
            "http_servers": has_path("/apps/http/servers"),
            "tls_automation": has_path("/apps/tls/automation"),
        }
        
        print("\n📊 系统状态:")
        for component, exists in status.items():
            icon = "✅" if exists else "❌"
            print(f"   {icon} {component}: {'已配置' if exists else '未配置'}")
        
        return status
    
    def cleanup_all(self, domains: List[str], base_domains: List[str] = None):
        """清理所有配置"""
        print("\n🧹 开始清理所有配置...")
        
        # 删除普通域名
        if domains:
            self.batch_delete(domains)
        
        # 删除通配符域名
        if base_domains:
            for base_domain in base_domains:
                wildcard_id = f"wildcard-{base_domain}"
                self.delete_domain(wildcard_id)

def demo_basic_operations():
    """演示基本操作"""
    print("=" * 60)
    print(" 演示1: 基本域名操作")
    print("=" * 60)
    
    dm = DomainManager()
    
    # 基本操作演示
    test_domains = {
        "api.example.com": "localhost:8080",
        "web.example.com": "localhost:3000",
        "admin.example.com": "localhost:9000"
    }
    
    # 检查初始状态
    dm.batch_check(list(test_domains.keys()))
    
    # 添加域名
    dm.batch_add(test_domains)
    
    # 验证添加结果
    dm.batch_check(list(test_domains.keys()))
    
    # 更新配置
    dm.update_domain("api.example.com", "localhost:8081")
    
    # 删除部分域名
    dm.batch_delete(["web.example.com"])
    
    # 最终检查
    dm.batch_check(list(test_domains.keys()))
    
    # 清理
    dm.cleanup_all(list(test_domains.keys()))

def demo_wildcard_operations():
    """演示通配符域名操作"""
    print("\n" + "=" * 60)
    print(" 演示2: 通配符域名操作")
    print("=" * 60)
    
    dm = DomainManager()
    
    base_domain = "dev.local"
    subdomains = {
        "app": "8090",
        "api": "8091",
        "admin": "8092",
        "monitoring": "8093"
    }
    
    # 添加通配符域名
    dm.add_wildcard_domain(base_domain)
    
    # 添加子域名
    print(f"\n📋 添加子域名到 {base_domain}...")
    for subdomain, port in subdomains.items():
        dm.add_subdomain(base_domain, subdomain, port)
    
    # 检查所有子域名状态
    subdomain_list = [f"{sub}.{base_domain}" for sub in subdomains.keys()]
    dm.batch_check(subdomain_list)
    
    # 删除部分子域名
    dm.batch_delete([f"monitoring.{base_domain}"])
    
    # 最终状态
    dm.batch_check(subdomain_list)
    
    # 清理
    dm.cleanup_all(subdomain_list, [base_domain])

def demo_advanced_scenarios():
    """演示高级使用场景"""
    print("\n" + "=" * 60)
    print(" 演示3: 高级使用场景")
    print("=" * 60)
    
    dm = DomainManager()
    
    # 场景1: 微服务架构域名管理
    microservices = {
        "auth.myapp.com": "localhost:8080",
        "user.myapp.com": "localhost:8081", 
        "order.myapp.com": "localhost:8082",
        "payment.myapp.com": "localhost:8083",
        "notification.myapp.com": "localhost:8084"
    }
    
    print("\n🏗️  场景1: 微服务架构域名管理")
    dm.batch_add(microservices)
    dm.get_system_status()
    
    # 场景2: 开发环境快速切换
    print("\n🔄 场景2: 开发环境快速切换")
    dev_configs = {
        "api.myapp.com": "localhost:3001",  # 切换到开发端口
        "web.myapp.com": "localhost:3000"   # 新增前端服务
    }
    dm.batch_add(dev_configs, force=True)  # 强制更新
    
    # 场景3: 生产环境部署前验证
    print("\n✅ 场景3: 部署前验证")
    all_domains = list(microservices.keys()) + list(dev_configs.keys())
    status = dm.batch_check(all_domains)
    
    configured_count = sum(1 for exists in status.values() if exists)
    print(f"\n📊 配置状态总结: {configured_count}/{len(all_domains)} 个域名已配置")
    
    # 清理
    dm.cleanup_all(all_domains)

def main():
    """主函数"""
    print("FastCaddy Python 域名管理示例")
    print("==============================")
    
    # 运行各种演示
    demo_basic_operations()
    demo_wildcard_operations() 
    demo_advanced_scenarios()
    
    print("\n" + "=" * 60)
    print(" 🎉 所有演示完成!")
    print("=" * 60)
    
    print("\n💡 域名管理最佳实践:")
    print("======================")
    print("1. 操作前总是检查域名状态 (has_id)")
    print("2. 批量操作使用专用方法提高效率") 
    print("3. 更新配置时使用 force=True 参数")
    print("4. 删除后验证操作结果")
    print("5. 使用通配符域名简化子域名管理")
    print("6. 定期检查系统状态确保配置正确")
    print("7. 开发和生产环境使用不同的域名前缀")
    print("8. 清理测试配置避免配置污染")

if __name__ == "__main__":
    main()