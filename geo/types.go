package geo

// ---------- High‑level JSON response ----------

type IPInfo struct {
	IP                string       `json:"ip"`                // queried IP
	Addr              string       `json:"addr,omitempty"`    // CIDR range
	Type              string       `json:"type,omitempty"`    // e.g. 基站/家庭宽带 …
	AS                *ASInfo      `json:"as,omitempty"`      // ASN details
	Country           *CountryInfo `json:"country,omitempty"` // geoname country
	RegisteredCountry *CountryInfo `json:"registered_country,omitempty"`
	Regions           []string     `json:"regions,omitempty"`       // 长名  [省, 市, 区]
	RegionsShort      []string     `json:"regions_short,omitempty"` // 简名 [川, 成都 …]
}

// ---------- Sub‑structures ----------

type ASInfo struct {
	Number uint   `json:"number"` // AS number
	Name   string `json:"name"`   // AS org (English)
	Info   string `json:"info"`   // Chinese ISP name
}

type CountryInfo struct {
	Code string `json:"code"` // CN / US …
	Name string `json:"name"` // 中国 / United States …
}

// ---------- Raw‑MMDB binding structs ----------

type CityRecord struct {
	Country struct {
		Iso   string            `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Iso   string            `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	Subdivisions []struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
}

type ASNRecord struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}

type CNRecord struct {
	Province string `maxminddb:"province"`
	City     string `maxminddb:"city"`
	District string `maxminddb:"districts"`
	ISP      string `maxminddb:"isp"`
	NetType  string `maxminddb:"net"`
}
