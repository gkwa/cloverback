package cloverback

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/patrickmn/go-cache"
)

var cacheRelPath = "cloverback/keys.db"

func Main(noExpunge bool) int {
	cachePath, err := getCachePath(cacheRelPath)
	if err != nil {
		panic(err)
	}
	mycache.LoadFile(cachePath)

	respCount := 1
	var pushSlice []Push
	pageCursor := ""

	for {
		req, err := genPushbulletRequest(pageCursor)
		if err != nil {
			return 1
		}

		respBody := requestPushbulletData(req)

		filePath := fmt.Sprintf("debug_resp_%02d.json", respCount)
		if err := savePushbulletResponseForDebug(respBody, filePath); err != nil {
			slog.Error("saving response body to file", "error", err)
			return 1
		}

		var pushBulletReply PushbulletHTTReply

		if err := json.Unmarshal(respBody, &pushBulletReply); err != nil { // Parse []byte to go struct pointer
			slog.Error("unmarshalling response body", "error", err)
			return 1
		}

		pushSlice = append(pushSlice, pushBulletReply.Pushes...)
		slog.Debug("push slice", "items", len(pushSlice))

		// // saveResponse(pushBulletReply)
		// if len(pushBulletReply.Pushes) > 0 {
		// 	buffer := genOrgMode(pushBulletReply)
		// 	writeBufferToClipboard(buffer)
		// 	// writeBufferToStdout(buffer)
		// }

		pageCursor = pushBulletReply.Cursor
		slog.Debug("cursor debug", "cursor", pageCursor)
		slog.Debug("pushbullet message", "count", len(pushBulletReply.Pushes))

		if len(pushBulletReply.Pushes) == 0 {
			break
		}
		respCount++
	}

	buffer := genOrgMode(pushSlice)
	writeBufferToClipboard(buffer)

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
		slog.Error("marhsalling response body", "error", err)
		return err
	}

	slog.Debug("caching", "key", cacheKey, "value", string(resultBytes))
	mycache.Set(cacheKey, string(resultBytes), cache.DefaultExpiration)
	return nil
}
