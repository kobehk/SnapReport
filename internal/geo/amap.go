package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AMapGeocoder 高德地图逆地理编码实现
type AMapGeocoder struct {
	APIKey string
	Client *http.Client
}

// NewAMapGeocoder 创建高德地图地理编码器
func NewAMapGeocoder(apiKey string) *AMapGeocoder {
	return &AMapGeocoder{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

// FlexibleString 处理高德地图返回的字段（可能是字符串或数组）
type FlexibleString string

func (c *FlexibleString) UnmarshalJSON(data []byte) error {
	// 尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*c = FlexibleString(str)
		return nil
	}

	// 尝试解析为字符串数组
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) > 0 {
			*c = FlexibleString(arr[0])
		} else {
			*c = FlexibleString("")
		}
		return nil
	}

	// 如果都不是，设为空字符串
	*c = FlexibleString("")
	return nil
}

// AMapReverseGeocodeResponse 高德地图逆地理编码响应
type AMapReverseGeocodeResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Regeocode struct {
		AddressComponent struct {
			Province     string         `json:"province"`
			City         FlexibleString `json:"city"`
			CityCode     string         `json:"citycode"`
			District     string         `json:"district"`
			Township     string         `json:"township"`
			Towncode     string         `json:"towncode"`
			Neighborhood struct {
				Name FlexibleString `json:"name"`
			} `json:"neighborhood"`
			Building struct {
				Name FlexibleString `json:"name"`
			} `json:"building"`
			StreetNumber struct {
				Street    FlexibleString `json:"street"`
				Number    FlexibleString `json:"number"`
				Direction FlexibleString `json:"direction"`
			} `json:"streetNumber"`
			BusinessAreas []struct {
				Name string `json:"name"`
			} `json:"businessAreas"`
		} `json:"addressComponent"`
		FormattedAddress string `json:"formatted_address"`
		Roads            []struct {
			Name      string  `json:"name"`
			Distance  float64 `json:"distance"`
			Direction string  `json:"direction"`
			Location  string  `json:"location"`
		} `json:"roads"`
	} `json:"regeocode"`
}

// ReverseGeocode 实现 Geocoder 接口
func (g *AMapGeocoder) ReverseGeocode(lat, lng float64) (string, string, string, error) {
	// 构造请求URL
	baseURL := "https://restapi.amap.com/v3/geocode/regeo"
	params := url.Values{}
	params.Set("key", g.APIKey)
	params.Set("location", fmt.Sprintf("%.6f,%.6f", lng, lat)) // 高德使用 "经度,纬度" 顺序
	params.Set("output", "json")
	params.Set("extensions", "base")
	params.Set("roadlevel", "1") // 返回道路信息

	reqURL := baseURL + "?" + params.Encode()
	resp, err := g.Client.Get(reqURL)
	if err != nil {
		return "", "", "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result AMapReverseGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", fmt.Errorf("parse response failed: %w", err)
	}

	if result.Status != "1" {
		return "", "", "", fmt.Errorf("amap api error: %s", result.Info)
	}

	// 提取城市信息
	cityStr := string(result.Regeocode.AddressComponent.City)
	if cityStr == "" {
		cityStr = result.Regeocode.AddressComponent.Province
	}

	// 提取道路信息
	road := ""
	if len(result.Regeocode.Roads) > 0 {
		road = result.Regeocode.Roads[0].Name
	}
	streetStr := string(result.Regeocode.AddressComponent.StreetNumber.Street)
	if road == "" && streetStr != "" {
		road = streetStr
	}

	// 高德不直接返回category，但我们可以根据道路类型推断
	category := "road"
	if len(result.Regeocode.Roads) > 0 {
		// 可以根据道路名称判断类型（如高速公路、国道等）
		roadName := road
		if containsAny(roadName, []string{"高速", "高速公路", "G", "国道"}) {
			category = "motorway"
		} else if containsAny(roadName, []string{"省道", "县道", "乡道"}) {
			category = "trunk"
		}
	}

	return cityStr, road, category, nil
}

func (g *AMapGeocoder) Provider() string {
	return "amap"
}
