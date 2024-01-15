package cloverback

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/patrickmn/go-cache"
)

var cacheRelPath = "cloverback/keys.db"

func Main(noExpunge bool) int {
	cachePath, err := getCachePath(cacheRelPath)
	if err != nil {
		panic(err)
	}
	mycache.LoadFile(cachePath)

	cacheDir := filepath.Dir(cachePath)

	respCount := 0
	var pushes []Push
	pageCursor := ""

	// results/pushes are paginated 20 per page
	for {
		req, err := genPushbulletRequest(pageCursor)
		if err != nil {
			slog.Error("pushbullet request generation", "error", err)
			return 1
		}

		respBody := requestPushbulletData(req)
		respCount++

		debugPath := filepath.Join(cacheDir, fmt.Sprintf("debug_resp_%02d.json", respCount))
		if err := savePushbulletResponseForDebug(respBody, debugPath); err != nil {
			slog.Error("saving response body to file", "error", err)
			return 1
		}

		var pushBulletReply PushbulletHTTReply
		if err := json.Unmarshal(respBody, &pushBulletReply); err != nil { // Parse []byte to go struct pointer
			slog.Error("unmarshalling response body", "error", err)
			return 1
		}

		slog.Debug("pushbullet message", "count", len(pushBulletReply.Pushes))

		pushes = append(pushes, pushBulletReply.Pushes...)
		slog.Debug("push slice", "items", len(pushes))

		if len(pushBulletReply.Pushes) == 0 {
			break
		}

		pageCursor = pushBulletReply.Cursor
	}

	backupPushbullets(pushes)
	buffer := genOrgMode(pushes, renderTmpl)

	clipboardWriter := &ClipboardWriter{}
	stdoutWriter := &StdoutWriter{}

	writeBuffer(clipboardWriter, &buffer)
	writeBuffer(stdoutWriter, &buffer)

	slog.Debug("caching", "cache path", cachePath)
	mycache.SaveFile(cachePath)
	slog.Debug("cache", "item count", mycache.ItemCount())
	getMostRecentCacheItem()
	saveValuesToFile()
	if !noExpunge {
		expungeAllPushbulletRecords()
	}

	slog.Info("pushes", "count", len(pushes))
	return 0
}

func backupPushbullets(result []Push) error {
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
