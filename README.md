mmdb 地址
https://github.com/P3TERX/GeoLite.mmdb

```bash
curl -O https://dl.google.com/go/go1.24.10.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.10.linux-amd64.tar.gz

alias go=/usr/local/go/bin/go
export GOPROXY=https://goproxy.cn,direct
do mod download
go build -o geoip main.go
cp -f geoip /usr/local/bin/
mkdir -p /usr/local/share/GeoIP
#cp -f GeoLite2-City.mmdb /usr/local/share/GeoIP/GeoLite2-City.mmdb
cp -f upgradegeoip /usr/local/bin
# 实时下载最新数据
upgradegeoip
```

### 测试

```bash
/usr/local/bin/geoip 92.118.39.58
--- IP Geo-Location 查询结果 ---
IP 地址: 92.118.39.58
数据库类型: GeoLite2-City (版本: 1764056431)
---------------------------------
国家/地区: 美国 (US)
城市: Dallas
经纬度: 32.7797, -96.8022
---------------------------------
```
