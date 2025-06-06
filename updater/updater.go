package updater

import (
	"fmt"
	"github.com/luolayo/geocn-go/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var Downloads = []struct {
	URL      string
	Filename string
}{
	{
		URL:      "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb",
		Filename: "GeoLite2-City.mmdb",
	},
	{
		URL:      "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb",
		Filename: "GeoLite2-ASN.mmdb",
	},
	{
		URL:      "https://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb",
		Filename: "GeoCN.mmdb",
	},
}

func StartDailyUpdater() {
	c := cron.New()
	// 每天0点执行
	_, err := c.AddFunc("0 0 * * *", updateAll)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// 立即执行一次，避免首次为空
	updateAll()

	c.Start()
}

func updateAll() {
	os.MkdirAll("data", 0755)
	logger.Info("Starting daily update of GeoLite databases...")

	for _, item := range Downloads {
		path := filepath.Join("data", item.Filename)
		logger.Info("Downloading GeoLite2-City.mmdb...")
		err := downloadFile(item.URL, path)
		if err != nil {
			logger.Error("Failed to download", zap.String("file", item.Filename), zap.Error(err))
		} else {
			logger.Info("Successfully downloaded", zap.String("file", item.Filename), zap.String("path", path))
		}
	}
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
func DownloadList() []struct{ URL, Filename string } {
	updateAll()
	return Downloads
}
