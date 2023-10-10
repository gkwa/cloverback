package cloverback

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/patrickmn/go-cache"
)

func getPushBulletAPIKey() string {
	apiKey := os.Getenv("PUSHBULLET_API_KEY")
	if apiKey == "" {
		slog.Error("PUSHBULLET_API_KEY environment variable is not set.")
		return ""
	}

	return apiKey
}

func Main() int {
	apiKey := getPushBulletAPIKey()
	resp := requestPushbulletData(apiKey)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response", "error", err.Error())
		return 1
	}

	cachePath, err := getCachePath(getCacheRelPath())
	if err != nil {
		panic(err)
	}

	var pushBulletReply PushbulletHTTReply

	mycache.LoadFile(cachePath)

	if err := json.Unmarshal(bodyBytes, &pushBulletReply); err != nil { // Parse []byte to go struct pointer
		slog.Error("unmarshalling response body", "error", err.Error())
		return 1
	}

	saveResponse(pushBulletReply)
	if len(pushBulletReply.Pushes) > 0 {
		buffer := genOrgMode(pushBulletReply)
		writeBufferToClipboard(buffer)
		writeBufferToStdout(buffer)
	}

	slog.Debug("caching", "cache path", cachePath)
	mycache.SaveFile(cachePath)
	slog.Debug("cache", "item count", mycache.ItemCount())
	getMostRecentCacheItem()
	saveValuesToFile()
	deleteAllPushbulletRecords()

	return 0
}

func saveResponse(result PushbulletHTTReply) error {
	resultBytes, err := json.Marshal(result)
	slog.Debug("json response", "json", resultBytes)
	if err != nil {
		slog.Error("marhsalling response body", "error", err.Error())
		return err
	}

	slog.Debug("caching", "key", cacheKey, "value", string(resultBytes))
	mycache.Set(cacheKey, string(resultBytes), cache.DefaultExpiration)
	return nil
}
