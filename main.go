// Program entry point.
//  1. Starts the daily updater (downloads MMDB files into ./data).
//  2. Loads the MMDBs with geo.MustOpen.
//  3. Exposes a single endpoint  /:ip  that returns JSON location info.
//
// Run with:  go run .         (listens :8080 by default)
//
//	go run . -http :9090
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/luolayo/geocn-go/geo"
	"github.com/luolayo/geocn-go/logger"
	"github.com/luolayo/geocn-go/updater"
	"github.com/oschwald/maxminddb-golang"
)

func main() {
	logger.Init("prod")
	logger.Info("Starting geocn-go server...")
	// start scheduled downloader: daily at 00:00 server time
	updater.StartDailyUpdater()

	listen := flag.String("http", ":8080", "HTTP listen address")
	flag.Parse()

	// open databases from ./data directory (created by updater)
	cityDB := geo.MustOpen("data/GeoLite2-City.mmdb")
	defer func(cityDB *maxminddb.Reader) {
		err := cityDB.Close()
		if err != nil {
			logger.S().Errorf("Failed to close city DB: %v", err)
		} else {
			logger.S().Info("City DB closed successfully")
		}
	}(cityDB)
	asnDB := geo.MustOpen("data/GeoLite2-ASN.mmdb")
	defer func(asnDB *maxminddb.Reader) {
		err := asnDB.Close()
		if err != nil {
			logger.S().Errorf("Failed to close ASN DB: %v", err)
		} else {
			logger.S().Info("ASN DB closed successfully")
		}
	}(asnDB)
	cnDB := geo.MustOpen("data/GeoCN.mmdb")
	defer func(cnDB *maxminddb.Reader) {
		err := cnDB.Close()
		if err != nil {
			logger.S().Errorf("Failed to close CN DB: %v", err)
		} else {
			logger.S().Info("CN DB closed successfully")
		}
	}(cnDB)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}
		c.Next()
	})
	r.GET("/:ip", func(c *gin.Context) {
		ip := c.Param("ip")
		res := geo.IPToAddress(ip, cityDB, asnDB, cnDB)
		c.JSON(200, res)
	})
	r.GET("/", func(c *gin.Context) {
		ip := c.ClientIP()
		res := geo.IPToAddress(ip, cityDB, asnDB, cnDB)
		c.JSON(200, res)
	})

	logger.S().Infof("Listening on %s", *listen)
	err := r.Run(*listen)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
