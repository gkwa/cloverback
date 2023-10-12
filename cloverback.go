package cloverback

import (
	"encoding/json"
	"log/slog"

	"github.com/patrickmn/go-cache"
)

var cacheRelPath = "cloverback/keys.db"

func Main(noExpunge bool) int {
	apiKey := getPushBulletAPIKey()
	bodyBytes := requestPushbulletData(apiKey)

	cachePath, err := getCachePath(cacheRelPath)
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
		slog.Debug("we got here")
		writeBufferToStdout(buffer)
		slog.Debug("we got here2")
	}

	slog.Debug("caching", "cache path", cachePath)
	mycache.SaveFile(cachePath)
	slog.Debug("cache", "item count", mycache.ItemCount())
	getMostRecentCacheItem()
	saveValuesToFile()
	if !noExpunge {
		expungeAllPushbulletRecords()
	}

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
