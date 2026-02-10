package model

type Report struct {
	ID        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Latitude  float64  `json:"lat"`
	Longitude float64  `json:"lng"`
	City      string   `json:"city"`
	RoadName  string   `json:"road_name"`
	IsHighway bool     `json:"is_highway"`
	VideoURL  string   `json:"video_url"`
	Status    string   `json:"status"`
	DeviceID  string   `json:"device_id"`
	Tags      []string `json:"tags"`
}
