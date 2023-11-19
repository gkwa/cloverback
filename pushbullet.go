package cloverback

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.pushbullet.com/v2/pushes"
)

func expungeAllPushbulletRecords() {
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

func savePushbulletResponseForDebug(bodyBytes []byte, filePath string) error {
	var prettyJSON PushbulletHTTReply
	if err := json.Unmarshal(bodyBytes, &prettyJSON); err != nil {
		return err
	}

	indentedJSON, err := json.MarshalIndent(prettyJSON, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, indentedJSON, 0o644); err != nil {
		return err
	}

	return nil
}

func getPushBulletAPIKey() string {
	apiKey := os.Getenv("PUSHBULLET_API_KEY")
	if apiKey == "" {
		slog.Error("PUSHBULLET_API_KEY environment variable is not set.")
	}

	return apiKey
}

func requestPushbulletData(req *http.Request) []byte {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("pushbullet", "type", "request", "status", "error", "error", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("pushbullet", "type", "response", "status", resp.StatusCode)
		return nil
	}
	slog.Debug("pushbullet", "type", "response", "status", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response", "error", err)
		return nil
	}

	return bodyBytes
}

func genPushbulletRequest(cursor string) (*http.Request, error) {
	apiKey := getPushBulletAPIKey()

	queryParams := map[string]string{
		"active":         "true",
		"modified_after": "1.4e+09",
	}

	if cursor != "" {
		queryParams["cursor"] = cursor
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		slog.Error("NewRequest", "error", err)
		return nil, err
	}

	req.Header.Set("Access-Token", apiKey)

	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	return req, nil
}
