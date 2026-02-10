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

func main() {
	// 1. Load Config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Dependencies
	memStore := store.NewMemoryStore()
	geocoder := geo.NewNominatimGeocoder(cfg.Nominatim.UserAgent)
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
