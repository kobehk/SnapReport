package api

import (
	"encoding/json"
	"net/http"

	"SnapReport/internal/service"
)

type Handler struct {
	Service *service.ReportService
}

func NewHandler(s *service.ReportService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/reports/prepare", h.prepare)
	mux.HandleFunc("/reports/send", h.send)
	mux.HandleFunc("/reports", h.list)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) prepare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		DeviceID    string   `json:"device_id"`
		Latitude    float64  `json:"lat"`
		Longitude   float64  `json:"lng"`
		DurationSec int      `json:"duration_sec"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if body.DurationSec <= 0 {
		body.DurationSec = 20
	}
	if body.DeviceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "device_id required"})
		return
	}

	req := service.PrepareRequest{
		DeviceID:    body.DeviceID,
		Latitude:    body.Latitude,
		Longitude:   body.Longitude,
		DurationSec: body.DurationSec,
		Tags:        body.Tags,
	}
	report, err := h.Service.Prepare(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	type response struct {
		ID        string  `json:"id"`
		Timestamp string  `json:"timestamp"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
		City      string  `json:"city"`
		RoadName  string  `json:"road_name"`
		IsHighway bool    `json:"is_highway"`
		VideoURL  string  `json:"video_url"`
		Status    string  `json:"status"`
		DeviceID  string  `json:"device_id"`
		Provider  string  `json:"provider"`
	}
	writeJSON(w, http.StatusOK, response{
		ID:        report.ID,
		Timestamp: report.Timestamp,
		Latitude:  report.Latitude,
		Longitude: report.Longitude,
		City:      report.City,
		RoadName:  report.RoadName,
		IsHighway: report.IsHighway,
		VideoURL:  report.VideoURL,
		Status:    report.Status,
		DeviceID:  report.DeviceID,
		Provider:  "nominatim",
	})
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	report, err := h.Service.Send(body.ID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":        report.ID,
		"status":    report.Status,
		"submitted": true,
	})
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Service.List())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
