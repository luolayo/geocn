# GeoCN-Go

🌐 **GeoCN-Go** is a lightweight IP-to-location REST API server written in Go. It combines [MaxMind GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) with [GeoCN](https://github.com/ljxi/GeoCN) to provide precise geographic and ISP data, including Chinese province, city, district, network type, and ASN information.

---

## ✨ Features

- ✅ IPv4 and IPv6 support
- ✅ Combines GeoLite2 and GeoCN.mmdb for enhanced accuracy in China
- ✅ Returns structured JSON data
- ✅ Fully dockerized with multi-architecture image (amd64 / arm64)
- ✅ Auto-downloads and updates mmdb data daily
- ✅ RESTful API interface

---

## 🚀 Quick Start

### 🐳 Option 1: Docker

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name geocn \
  luola/geocn-go:latest
```

### 🐙 Option 2: Docker Compose

Make sure you have a `docker-compose.yml` like this:

```yaml
services:
  geocn:
    image: luola/geocn-go:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

Then run:

```bash
docker compose up -d
```

---

## 📡 API Example

Query:

```bash
GET http://localhost:8080/114.114.114.114
```

Response:

```json
{
  "ip": "114.114.114.114",
  "addr": "114.114.114.0/24",
  "type": "家庭宽带",
  "as": {
    "number": 4134,
    "name": "Chinanet",
    "info": "中国电信"
  },
  "country": {
    "code": "CN",
    "name": "中国"
  },
  "registered_country": {
    "code": "CN",
    "name": "中国"
  },
  "regions": ["江苏省", "南京市"],
  "regions_short": ["江苏", "南京"]
}
```

---

## 📦 MMDB Sources

| File               | Description                    | Source                                                                           |
|--------------------|--------------------------------|----------------------------------------------------------------------------------|
| GeoLite2-City.mmdb | Country, city, region data     | [MaxMind GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) |
| GeoLite2-ASN.mmdb  | ASN & ISP info                 | MaxMind ASN DB                                                                   |
| GeoCN.mmdb         | Chinese ISP, region, type data | [ljxi/GeoCN](https://github.com/ljxi/GeoCN)                                      |

These files are auto-downloaded daily and stored in the `./data/` directory.

---

## 🏗️ Local Development (Optional)

```bash
go build -o geocn main.go
./geocn
```

---

## 🐳 Docker Multi-Arch Support

| Platform       | Supported |
|----------------|-----------|
| `linux/amd64`  | ✅ Yes     |
| `linux/arm64`  | ✅ Yes     |

Docker image: [luola/geocn-go on Docker Hub](https://hub.docker.com/r/luola/geocn-go)

---

## 📜 License

MIT License.  
For educational and research use only. Data from MaxMind and GeoCN.