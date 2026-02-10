package service

import (
	"fmt"
	"strconv"
	"time"

	"SnapReport/internal/ddpai"
	"SnapReport/internal/geo"
	"SnapReport/internal/model"
	"SnapReport/internal/store"
)

type ReportService struct {
	Store    store.Store
	Geocoder geo.Geocoder
	DDPai    *ddpai.Client
}

func NewReportService(s store.Store, g geo.Geocoder, d *ddpai.Client) *ReportService {
	return &ReportService{
		Store:    s,
		Geocoder: g,
		DDPai:    d,
	}
}

type PrepareRequest struct {
	DeviceID    string
	Latitude    float64
	Longitude   float64
	DurationSec int
	Tags        []string
}

func (s *ReportService) Prepare(req PrepareRequest) (*model.Report, error) {
	city, road, category, err := s.Geocoder.ReverseGeocode(req.Latitude, req.Longitude)
	if err != nil {
		// Log error but continue, don't fail the whole request
		fmt.Printf("Warning: geocode failed: %v\n", err)
		city = "Unknown"
		road = "Unknown"
	}

	isHighway := geo.ClassifyHighway(category, road)

	videoURL, err := s.DDPai.CaptureRecentVideo(req.DeviceID, req.DurationSec)
	if err != nil {
		return nil, fmt.Errorf("capture video failed: %w", err)
	}

	id := s.newID()
	now := time.Now().UTC().Format(time.RFC3339)
	report := model.Report{
		ID:        id,
		Timestamp: now,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		City:      city,
		RoadName:  road,
		IsHighway: isHighway,
		VideoURL:  videoURL,
		Status:    "prepared",
		DeviceID:  req.DeviceID,
		Tags:      req.Tags,
	}
	s.Store.Save(report)
	return &report, nil
}

func (s *ReportService) Send(id string) (*model.Report, error) {
	report, ok := s.Store.Get(id)
	if !ok {
		return nil, fmt.Errorf("report not found")
	}
	report.Status = "submitted"
	s.Store.Save(report)
	return &report, nil
}

func (s *ReportService) List() []model.Report {
	return s.Store.List()
}

func (s *ReportService) newID() string {
	now := time.Now().UTC().UnixNano()
	return "rep_" + strconv.FormatInt(now, 36)
}
