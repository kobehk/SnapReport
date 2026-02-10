package main

import (
	"log"
	"net/http"
	"strconv"

	"SnapReport/internal/api"
	"SnapReport/internal/config"
	"SnapReport/internal/ddpai"
	"SnapReport/internal/geo"
	"SnapReport/internal/service"
	"SnapReport/internal/store"
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

	// 5. Setup Router
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// 6. Start Server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("SnapReport backend listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
