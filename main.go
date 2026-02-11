package main

import (
	"fmt"
	"log"

	"SnapReport/internal/api"
	"SnapReport/internal/config"
	"SnapReport/internal/ddpai"
	"SnapReport/internal/geo"
	"SnapReport/internal/service"
	"SnapReport/internal/store"

	"github.com/gin-gonic/gin"
)

// maskAPIKey 隐藏API密钥的中间部分，仅显示首尾字符
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

func main() {
	// 1. Load Config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Dependencies
	memStore := store.NewMemoryStore()

	// 根据配置选择地理编码器
	var geocoder geo.Geocoder
	switch cfg.Geocoder.Type {
	case "amap":
		if cfg.Geocoder.APIKey == "" {
			log.Fatal("AMap geocoder requires API key. Please set geocoder.api_key in config.yaml")
		}
		geocoder = geo.NewAMapGeocoder(cfg.Geocoder.APIKey)
		log.Printf("Using AMap geocoder with API key: %s", maskAPIKey(cfg.Geocoder.APIKey))
	default: // "nominatim" 或未指定
		geocoder = geo.NewNominatimGeocoder(cfg.Geocoder.UserAgent)
		log.Printf("Using Nominatim geocoder with user agent: %s", cfg.Geocoder.UserAgent)
	}

	ddpaiClient := ddpai.NewClient(
		cfg.DDPai.BaseURL,
		cfg.DDPai.TimeoutSeconds,
		cfg.DDPai.MockMode,
	)

	// 3. Initialize Service
	svc := service.NewReportService(memStore, geocoder, ddpaiClient)

	// 4. Initialize Handler
	handler := api.NewHandler(svc)

	// 5. Setup Gin Router
	router := gin.Default()
	handler.RegisterGinRoutes(router)

	// 6. Start Server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("SnapReport backend listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
