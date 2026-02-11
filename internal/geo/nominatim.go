package geo

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
)

type Geocoder interface {
	ReverseGeocode(lat, lng float64) (city, road, category string, err error)
	Provider() string
}

type NominatimGeocoder struct {
	UserAgent string
	Client    *http.Client
}

func NewNominatimGeocoder(userAgent string) *NominatimGeocoder {
	// 创建自定义的HTTP客户端，不使用代理，并配置TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // 要求TLS 1.2或更高
	}
	transport := &http.Transport{
		Proxy:           nil, // 禁用代理
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
		// 防止自动重定向到HTTP，保持HTTPS连接
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 如果重定向到HTTP，返回错误以保持HTTPS
			if req.URL.Scheme == "http" {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	return &NominatimGeocoder{
		UserAgent: userAgent,
		Client:    client,
	}
}

func (g *NominatimGeocoder) ReverseGeocode(lat, lng float64) (string, string, string, error) {
	type nominatimResp struct {
		Address struct {
			City          string `json:"city"`
			County        string `json:"county"`
			State         string `json:"state"`
			Country       string `json:"country"`
			Road          string `json:"road"`
			Neighbourhood string `json:"neighbourhood"`
			Suburb        string `json:"suburb"`
		} `json:"address"`
		Category string `json:"category"`
		Type     string `json:"type"`
	}
	url := "https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=" +
		strconv.FormatFloat(lat, 'f', 6, 64) + "&lon=" + strconv.FormatFloat(lng, 'f', 6, 64)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", g.UserAgent)
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	var nr nominatimResp
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		return "", "", "", err
	}
	city := nr.Address.City
	if city == "" {
		if nr.Address.County != "" {
			city = nr.Address.County
		} else if nr.Address.State != "" {
			city = nr.Address.State
		}
	}
	road := nr.Address.Road
	category := nr.Type
	if category == "" {
		category = nr.Category
	}
	return city, road, category, nil
}

func (g *NominatimGeocoder) Provider() string {
	return "nominatim"
}
