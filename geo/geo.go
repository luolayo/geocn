package geo

import (
	"net"
	"sync/atomic"

	"github.com/oschwald/maxminddb-golang"
)

// Minimal ASN → Chinese ISP mapping.  Add more entries as needed.
var asnMap = map[uint]string{
	9812:  "东方有线",
	9389:  "中国长城",
	17962: "天威视讯",
	17429: "歌华有线",
	7497:  "科技网",
	24139: "华数",
	9801:  "中关村",
	4538:  "教育网",
	24151: "CNNIC",

	38019: "中国移动", 139080: "中国移动", 9808: "中国移动", 24400: "中国移动", 134810: "中国移动", 24547: "中国移动",
	56040: "中国移动", 56041: "中国移动", 56042: "中国移动", 56044: "中国移动", 132525: "中国移动", 56046: "中国移动",
	56047: "中国移动", 56048: "中国移动", 59257: "中国移动", 24444: "中国移动",
	24445: "中国移动", 137872: "中国移动", 9231: "中国移动", 58453: "中国移动",

	4134: "中国电信", 4812: "中国电信", 23724: "中国电信", 136188: "中国电信", 137693: "中国电信", 17638: "中国电信",
	140553: "中国电信", 140061: "中国电信", 136195: "中国电信", 17799: "中国电信", 139018: "中国电信",
	133776: "中国电信", 58772: "中国电信", 146966: "中国电信", 63527: "中国电信", 58539: "中国电信", 58540: "中国电信",
	141998: "中国电信", 138169: "中国电信", 139203: "中国电信", 58563: "中国电信", 137690: "中国电信", 63838: "中国电信",
	137694: "中国电信", 137698: "中国电信", 136167: "中国电信", 148969: "中国电信", 134764: "中国电信",
	134770: "中国电信", 148981: "中国电信", 134774: "中国电信", 136190: "中国电信", 140647: "中国电信",
	132225: "中国电信", 140485: "中国电信", 4811: "中国电信", 131285: "中国电信", 137689: "中国电信",
	137692: "中国电信", 140636: "中国电信", 140638: "中国电信", 140345: "中国电信", 38283: "中国电信",
	140292: "中国电信", 140903: "中国电信", 17897: "中国电信", 134762: "中国电信", 139019: "中国电信",
	141739: "中国电信", 141771: "中国电信", 134419: "中国电信", 140276: "中国电信", 58542: "中国电信",
	140278: "中国电信", 139767: "中国电信", 137688: "中国电信", 137691: "中国电信", 4809: "中国电信",
	58466: "中国电信", 137687: "中国电信", 134756: "中国电信", 134760: "中国电信",
	133774: "中国电信", 133775: "中国电信", 4816: "中国电信", 134768: "中国电信",
	58461: "中国电信", 58520: "中国电信", 131325: "中国电信",

	4837: "中国联通", 4808: "中国联通", 134542: "中国联通", 134543: "中国联通", 10099: "中国联通",
	140979: "中国联通", 138421: "中国联通", 17621: "中国联通", 17622: "中国联通", 17816: "中国联通",
	140726: "中国联通", 17623: "中国联通", 136958: "中国联通", 9929: "中国联通",
	140716: "中国联通", 4847: "中国联通", 136959: "中国联通", 135061: "中国联通", 139007: "中国联通",

	59019:  "金山云",
	135377: "优刻云",
	45062:  "网易云",
	137718: "火山引擎",
	37963:  "阿里云", 45102: "阿里云国际",
	45090: "腾讯云", 132203: "腾讯云国际",
	55967: "百度云", 38365: "百度云",
	58519: "华为云", 55990: "华为云", 136907: "华为云",
	4609:   "澳門電訊",
	134773: "珠江宽频",
	1659:   "台湾教育网",
	8075:   "微软云",
	17421:  "中华电信",
	3462:   "HiNet",
	13335:  "Cloudflare",
	55960:  "亚马逊云", 14618: "亚马逊云", 16509: "亚马逊云",
	15169: "谷歌云", 396982: "谷歌云", 36492: "谷歌云",
}

var (
	cityDB atomic.Value // 线程安全读写
	asnDB  atomic.Value
	cnDB   atomic.Value
)

// IPToAddress converts an IP string to structured info.
// It mirrors the original Python behaviour but in Go.
func IPToAddress(ip string, cityDB, asnDB, cnDB *maxminddb.Reader) *IPInfo {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return &IPInfo{IP: ip}
	}

	// ---- City lookup -------------------------------------------------
	var cityRec CityRecord
	cityNet, cityOK, _ := cityDB.LookupNetwork(parsed, &cityRec)

	// ---- CN (Chinese enhanced) lookup --------------------------------
	var cnRec CNRecord
	cnNet, cnOK, _ := cnDB.LookupNetwork(parsed, &cnRec)

	// ---- ASN lookup --------------------------------------------------
	var asnRec ASNRecord
	_ = asnDB.Lookup(parsed, &asnRec)

	info := &IPInfo{IP: ip}

	// choose CIDR: prefer CN db if hit
	switch {
	case cnOK:
		info.Addr = cnNet.String()
	case cityOK:
		info.Addr = cityNet.String()
	}

	// basic country fields
	info.Country = &CountryInfo{Code: cityRec.Country.Iso, Name: pickName(cityRec.Country.Names)}
	info.RegisteredCountry = &CountryInfo{Code: cityRec.RegisteredCountry.Iso, Name: pickName(cityRec.RegisteredCountry.Names)}

	// ASN information
	if asnRec.AutonomousSystemNumber != 0 {
		info.AS = &ASInfo{
			Number: asnRec.AutonomousSystemNumber,
			Name:   asnRec.AutonomousSystemOrganization,
			Info:   asnMap[asnRec.AutonomousSystemNumber],
		}
	}
	// override/patch with CN ISP if present
	if cnRec.ISP != "" {
		if info.AS == nil {
			info.AS = &ASInfo{}
		}
		if info.AS.Info == "" {
			info.AS.Info = cnRec.ISP
		}
	}
	if cnRec.NetType != "" {
		info.Type = cnRec.NetType
	}

	// Compose regions then apply simplification rule
	raw := []string{
		cnRec.Province,
		cnRec.City,
		cnRec.District,
		func() string {
			if len(cityRec.Subdivisions) > 0 {
				return pickName(cityRec.Subdivisions[0].Names)
			}
			return ""
		}(),
		pickName(cityRec.City.Names),
	}
	info.Regions = simplifyRegions(raw)

	// Build regions_short based on final regions
	for _, r := range info.Regions {
		info.RegionsShort = append(info.RegionsShort, shortName(r))
	}
	return info
}
