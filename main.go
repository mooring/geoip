package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings" // 用于错误信息处理
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

// GeoRecord 定义结构体来映射 GeoLite2-City 数据库中的字段
// 使用 interface{} 提高解码兼容性
type GeoRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
		// 更改为 interface{} 以最大化兼容性
		Name map[string]interface{} `maxminddb:"names"`
	} `maxminddb:"country"`

	City struct {
		// 更改为 interface{} 以最大化兼容性
		Name map[string]interface{} `maxminddb:"names"`
	} `maxminddb:"city"`

	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
}

// 辅助函数：安全地从 map[string]interface{} 中提取 string
func getMapName(names map[string]interface{}, lang string) string {
	if names == nil {
		return ""
	}
	if val, ok := names[lang]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func main() {
	// --- 配置 ---
	const dbPath = "/usr/local/share/GeoIP/GeoLite2-City.mmdb"
	// -----------------

	// 1. 检查命令行参数
	if len(os.Args) != 2 {
		fmt.Printf("用法: %s <IP地址>\n", os.Args[0])
		fmt.Println("示例: geoip 2.228.129.230")
		os.Exit(1)
	}

	ipStr := os.Args[1]

	// 使用 netip.ParseAddr 解析 IP 字符串
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		log.Fatalf("错误: IP地址格式无效或不支持: %s. %v", ipStr, err)
	}

	// 2. 打开 GeoIP 数据库文件
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Fatalf("无法打开数据库文件 %s: %v", dbPath, err)
	}
	defer db.Close()

	// 3. 查询 IP 地址
	var record GeoRecord
	// 适配本地库签名：只接收一个返回值
	result := db.Lookup(addr)

	// 4. 检查查询结果和解码
	// 关键步骤：手动解码到结构体，并将 Decode() 的返回值用于错误处理
	err = result.Decode(&record)
	if err != nil {
		errStr := err.Error()

		// 检查是否是 "找不到记录" 的错误
		if strings.Contains(errStr, "Cannot locate") || strings.Contains(errStr, "not found") {
			log.Fatalf("查询 IP %s 失败: 在数据库中找不到该 IP 的记录。", ipStr)
		}

		// 其他解码错误
		log.Fatalf("查询 IP %s 成功，但解码数据到结构体失败: %s", ipStr, errStr)
	}

	// 5. 打印查询结果
	fmt.Println("--- IP GEO 查询结果 ---")
	now := time.Unix(int64(db.Metadata.BuildEpoch), 0)
	ver := now.Format("20060102")
	fmt.Printf("查询地址: %s\n", ipStr)
	fmt.Printf("数据版本: %s(%s)\n", db.Metadata.DatabaseType, ver)
	fmt.Println("---------------------------------")

	// 优先显示中文名称，否则回退到英文
	countryName := getMapName(record.Country.Name, "zh-CN")
	if countryName == "" {
		countryName = getMapName(record.Country.Name, "en")
	}
	cityName := getMapName(record.City.Name, "zh-CN")
	if cityName == "" {
		cityName = getMapName(record.City.Name, "en")
	}

	// 输出详细信息
	fmt.Printf("国家地区: %s (%s)\n", countryName, record.Country.ISOCode)
	fmt.Printf("所属城市: %s\n", cityName)
	fmt.Printf("经 纬 度: %.4f, %.4f\n", record.Location.Latitude, record.Location.Longitude)

	fmt.Println("---------------------------------")
}
