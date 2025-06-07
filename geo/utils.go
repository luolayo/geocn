package geo

import (
	"log"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

// MustOpen opens an mmdb file or terminates the program.
func MustOpen(path string) *maxminddb.Reader {
	db, err := maxminddb.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	return db
}

// pickName returns zh‑CN name when available, otherwise any available name.
func pickName(names map[string]string) string {
	if v := names["zh-CN"]; v != "" {
		return v
	}
	for _, v := range names { // first map entry
		return v
	}
	return ""
}

// shortName trims Chinese suffixes (省、市、区…).
func shortName(name string) string {
	suf := []string{"省", "市", "区", "壮族自治区", "回族自治区", "维吾尔自治区", "自治区", "特别行政区"}
	for _, s := range suf {
		if strings.HasSuffix(name, s) {
			return strings.TrimSuffix(name, s)
		}
	}
	return name
}

// simplifyRegions:
//   - removes blanks
//   - if resulting slice has >2 elems, drops the last two items
//   - returns slice as‑is when len<=2
func simplifyRegions(list []string) []string {
	out := make([]string, 0, len(list))
	for _, v := range list {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) > 2 {
		return out[:len(out)-2]
	}
	return out
}

// Ensure data directory exists at program start (used by updater tests)
func ensureDataDir() {
	_ = os.MkdirAll("data", 0o755)
}

func ReloadDatabases() {
	newCityDB, _ := maxminddb.Open("data/GeoLite2-City.mmdb")
	newASNDB, _ := maxminddb.Open("data/GeoLite2-ASN.mmdb")
	newCNDB, _ := maxminddb.Open("data/GeoCN.mmdb")

	// 关闭旧 reader
	if old := cityDB.Swap(newCityDB); old != nil {
		err := old.(*maxminddb.Reader).Close()
		if err != nil {
			return
		}
	}
	if old := asnDB.Swap(newASNDB); old != nil {
		err := old.(*maxminddb.Reader).Close()
		if err != nil {
			return
		}
	}
	if old := cnDB.Swap(newCNDB); old != nil {
		old.(*maxminddb.Reader).Close()
	}
}

func GetReaders() (*maxminddb.Reader, *maxminddb.Reader, *maxminddb.Reader) {
	return cityDB.Load().(*maxminddb.Reader), asnDB.Load().(*maxminddb.Reader), cnDB.Load().(*maxminddb.Reader)
}
