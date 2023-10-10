package cloverback

import (
	"log/slog"
	"net/http"
)

var apiURL string

func init() {
	apiURL = "https://api.pushbullet.com/v2/pushes"
}

func requestPushbulletData(apiKey string) *http.Response {
	client := &http.Client{}

	queryParams := map[string]string{
		"active":         "true",
		"modified_after": "1.4e+09",
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		slog.Error("NewRequest", "error", err.Error())
		return nil
	}

	req.Header.Set("Access-Token", apiKey)

	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("pushbullet", "type", "request", "status", "error", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("pushbullet", "type", "response", "status", resp.StatusCode)
		return nil
	}
	slog.Debug("pushbullet", "type", "response", "status", resp.StatusCode)

	return resp
}

func deleteAllPushbulletRecords() {
	apiKey := getPushBulletAPIKey()

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		slog.Debug("Error creating request:", err)
		return
	}

	req.Header.Set("Access-Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		slog.Debug("pushbullet", "type", "response", "status", resp.StatusCode)
	} else {
		slog.Error("pushbullet", "type", "response", "status", resp.StatusCode)
	}
}
