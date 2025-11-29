mmdb 地址
https://github.com/P3TERX/GeoLite.mmdb

```bash
export GOPROXY=https://goproxy.cn,direct
do mod download
go build -o geoip main.go
cp geoip /usr/local/bin/
```

### 测试

```bash
./geoip 92.118.39.58
--- IP Geo-Location 查询结果 ---
IP 地址: 92.118.39.58
数据库类型: GeoLite2-City (版本: 1764056431)
---------------------------------
国家/地区: 美国 (US)
城市: Dallas
经纬度: 32.7797, -96.8022
---------------------------------
```
