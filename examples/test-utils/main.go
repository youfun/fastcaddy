package main

import (
	"fmt"

	"github.com/youfun/fastcaddy/internal/config"
	"github.com/youfun/fastcaddy/internal/utils"
)

// 测试工具函数的正确性
func main() {
	fmt.Println("FastCaddy Go 版本 - 工具函数测试")
	fmt.Println("================================")

	// 测试路径处理函数
	fmt.Println("\n1. 测试路径处理函数:")
	testPaths := []string{
		"/apps/http/servers",
		"apps/http/servers/",
		"apps/http/servers",
		"/",
		"",
	}

	for _, path := range testPaths {
		normalized := utils.NormalizePath(path)
		cleaned := utils.CleanPath(path)
		fmt.Printf("输入: '%s' -> 标准化: '%s', 清理: '%s'\n", path, normalized, cleaned)
	}

	// 测试路径分割和连接
	fmt.Println("\n2. 测试路径分割和连接:")
	testPath := "/apps/http/servers/srv0/routes"
	components := utils.SplitPath(testPath)
	rejoined := utils.JoinPath(components...)
	fmt.Printf("原路径: %s\n", testPath)
	fmt.Printf("分割后: %v\n", components)
	fmt.Printf("重新连接: %s\n", rejoined)

	// 测试配置管理函数
	fmt.Println("\n3. 测试配置管理函数:")
	
	// 测试 path2keys 和 keys2path
	path := "/apps/tls/automation/policies"
	keys := config.PathToKeys(path)
	reconstructed := config.KeysToPath(keys...)
	fmt.Printf("原路径: %s\n", path)
	fmt.Printf("转换为键: %v\n", keys)
	fmt.Printf("重构路径: %s\n", reconstructed)

	// 测试嵌套字典操作
	fmt.Println("\n4. 测试嵌套字典操作:")
	dict := make(map[string]interface{})
	config.NestedSetDict(dict, "test_value", "level1", "level2", "level3")
	fmt.Printf("设置嵌套值后: %+v\n", dict)

	// 测试主机名和 URL 验证
	fmt.Println("\n5. 测试验证函数:")
	testHosts := []string{"api.example.com", "localhost", "invalid host", ""}
	testURLs := []string{"localhost:8080", "127.0.0.1:3000", "invalid", ""}

	for _, host := range testHosts {
		valid := utils.ValidateHost(host)
		fmt.Printf("主机 '%s': %t\n", host, valid)
	}

	for _, url := range testURLs {
		valid := utils.ValidateURL(url)
		fmt.Printf("URL '%s': %t\n", url, valid)
	}

	// 测试字符串工具函数
	fmt.Println("\n6. 测试字符串工具函数:")
	slice := []string{"apple", "banana", "cherry"}
	fmt.Printf("切片 %v 包含 'banana': %t\n", slice, utils.StringSliceContains(slice, "banana"))
	fmt.Printf("切片 %v 包含 'grape': %t\n", slice, utils.StringSliceContains(slice, "grape"))

	defaultValue := utils.DefaultIfEmpty("", "default")
	nonEmptyValue := utils.DefaultIfEmpty("value", "default")
	fmt.Printf("空值默认处理: '%s'\n", defaultValue)
	fmt.Printf("非空值默认处理: '%s'\n", nonEmptyValue)

	// 测试映射合并
	map1 := map[string]string{"a": "1", "b": "2"}
	map2 := map[string]string{"b": "3", "c": "4"}
	merged := utils.MergeStringMaps(map1, map2)
	fmt.Printf("合并映射 %v + %v = %v\n", map1, map2, merged)

	fmt.Println("\n✅ 所有工具函数测试完成!")
	
	// 验证核心函数是否按预期工作
	fmt.Println("\n7. 验证核心功能一致性:")
	
	// 验证 PathToKeys 和 KeysToPath 的往返转换
	originalPath := "/apps/http/servers/srv0"
	keys1 := config.PathToKeys(originalPath)
	reconstructedPath := config.KeysToPath(keys1...)
	
	if reconstructedPath == originalPath {
		fmt.Printf("✅ 路径往返转换正确: %s\n", originalPath)
	} else {
		fmt.Printf("❌ 路径往返转换错误: %s -> %s\n", originalPath, reconstructedPath)
	}
	
	// 验证嵌套字典设置的正确性
	testDict := make(map[string]interface{})
	config.NestedSetDict(testDict, "value", "a", "b", "c")
	
	// 检查是否正确设置了嵌套值
	if val, ok := testDict["a"].(map[string]interface{})["b"].(map[string]interface{})["c"]; ok && val == "value" {
		fmt.Println("✅ 嵌套字典设置正确")
	} else {
		fmt.Println("❌ 嵌套字典设置错误")
	}
	
	// 验证类型一致性
	fmt.Println("\n8. Python 函数对应关系验证:")
	pythonGoMapping := map[string]string{
		"get_id()":            "client.GetIDURL()",
		"get_path()":          "client.GetConfigURL()",
		"gid()":               "client.GetByID()",
		"gcfg()":              "client.GetConfig()",
		"has_id()":            "client.HasID()",
		"has_path()":          "client.HasPath()",
		"pid()":               "client.PutByID()",
		"pcfg()":              "client.PutConfig()",
		"del_id()":            "client.DeleteByID()",
		"nested_setdict()":    "config.NestedSetDict()",
		"path2keys()":         "config.PathToKeys()",
		"keys2path()":         "config.KeysToPath()",
		"init_path()":         "manager.InitPath()",
		"add_reverse_proxy()": "routes.AddReverseProxy()",
		"add_wildcard_route()": "routes.AddWildcardRoute()",
		"setup_caddy()":       "fastcaddy.SetupCaddy()",
	}
	
	fmt.Println("Python 函数 -> Go 函数映射:")
	for py, go_ := range pythonGoMapping {
		fmt.Printf("  %s -> %s\n", py, go_)
	}
}