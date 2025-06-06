package test

import (
	"github.com/luolayo/geocn-go/geo"
	"math/rand"
	"testing"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

var iplist = []string{
	"47.106.144.184",
	"221.224.25.37",
	"60.204.145.212",
	"218.6.120.111",
	"139.196.111.167",
	"118.31.1.154",
	"183.60.141.17",
	"114.231.82.200",
	"183.236.232.160",
	"121.5.130.51",
	"240e:476::",
}

func TestIPConversion200Each(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// 打开 mmdb
	cityDB, err := maxminddb.Open("../data/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatalf("open city db: %v", err)
	}
	defer cityDB.Close()
	asnDB, _ := maxminddb.Open("../data/GeoLite2-ASN.mmdb")
	defer asnDB.Close()
	cnDB, _ := maxminddb.Open("../data/GeoCN.mmdb")
	defer cnDB.Close()

	for i := 0; i < len(iplist); i++ {
		ip := iplist[i]
		info := geo.IPToAddress(ip, cityDB, asnDB, cnDB)
		t.Logf("ipv4 %s regions short too many: %v", ip, info.RegionsShort)
		if info.IP != ip {
			t.Fatalf("ipv4 mismatch: %s vs %s", ip, info.IP)
		}
		if info.Addr == "" {
			t.Fatalf("ipv4 %s addr empty", ip)
		}
	}

}
