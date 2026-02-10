package ddpai

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	BaseURL  string
	Client   *http.Client
	MockMode bool
}

func NewClient(baseURL string, timeoutSeconds int, mockMode bool) *Client {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 5
	}
	return &Client{
		BaseURL:  baseURL,
		Client:   &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
		MockMode: mockMode,
	}
}

func (c *Client) CaptureRecentVideo(deviceID string, durationSec int) (string, error) {
	session, err := c.getSession()
	if err != nil {
		if c.MockMode {
			return c.mockURL(deviceID, durationSec), nil
		}
		// Try to proceed even if session fails, some firmwares might not need it?
		// Or just fail. Let's fail if not mock mode, consistent with logic.
		// Actually original logic was a bit mixed. Let's stick to fail-fast unless mock.
		return "", err
	}

	_ = c.setSuperDownload(session, true)
	list, err := c.getPlaybackList(session)
	if err != nil || len(list) == 0 {
		if c.MockMode {
			return c.mockURL(deviceID, durationSec), nil
		}
		if err != nil {
			return "", err
		}
		return "", nil // Or error "no video found"
	}

	item := list[len(list)-1]
	name := ""
	if v, ok := item["name"].(string); ok && v != "" {
		name = v
	} else if v, ok := item["file"].(string); ok && v != "" {
		name = v
	}
	if name == "" {
		if c.MockMode {
			return c.mockURL(deviceID, durationSec), nil
		}
		return "", nil
	}

	url := c.BaseURL + "/cmd.cgi?cmd=API_FileDownloadReq&file=" + name
	if session != "" {
		url += "&session=" + session
	}
	return url, nil
}

func (c *Client) mockURL(deviceID string, durationSec int) string {
	return "ddpai://device/" + deviceID + "/clip?duration=" + strconv.Itoa(durationSec)
}

func (c *Client) getSession() (string, error) {
	u := c.BaseURL + "/cmd.cgi?cmd=API_SessionReq"
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var m map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if v, ok := m["session"].(string); ok {
		return v, nil
	}
	if v, ok := m["sid"].(string); ok {
		return v, nil
	}
	return "", nil
}

func (c *Client) setSuperDownload(session string, enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	u := c.BaseURL + "/cmd.cgi?cmd=API_SuperDownloadReq&enable=" + val
	if session != "" {
		u += "&session=" + session
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) getPlaybackList(session string) ([]map[string]any, error) {
	u := c.BaseURL + "/cmd.cgi?cmd=API_PlaybackListReq"
	if session != "" {
		u += "&session=" + session
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode into interface{} to handle both array and object
	var raw any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if arr, ok := raw.([]any); ok {
		out := make([]map[string]any, 0, len(arr))
		for _, it := range arr {
			if m, ok := it.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out, nil
	}

	if obj, ok := raw.(map[string]any); ok {
		if v, ok := obj["list"].([]any); ok {
			out := make([]map[string]any, 0, len(v))
			for _, it := range v {
				if m, ok := it.(map[string]any); ok {
					out = append(out, m)
				}
			}
			return out, nil
		}
	}

	return nil, nil
}
