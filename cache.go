package cloverback

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/adrg/xdg"

	"github.com/patrickmn/go-cache"
)

var (
	mycache   *cache.Cache
	timestamp int64
	cacheKey  string
)

func init() {
	mycache = cache.New(12*time.Hour, 24*time.Hour)
	timestamp = time.Now().Unix()
	cacheKey = fmt.Sprintf("cloverback_%d", timestamp)
}

func getCachePath(configRelPath string) (string, error) {
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}

	dirPerm := os.FileMode(0o700)

	d := filepath.Dir(configFilePath)

	if err := os.MkdirAll(d, dirPerm); err != nil {
		return "", err
	}

	slog.Debug("cache", "path", configFilePath)
	logPathStats(configFilePath)
	return configFilePath, nil
}

func logPathStats(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("stat", "path", filePath, "error", err)
		return
	}

	fileUID := fileInfo.Sys().(*syscall.Stat_t).Uid

	// Use the user package to get the user information
	u, err := user.LookupId(fmt.Sprintf("%d", fileUID))
	if err != nil {
		slog.Error("user info", "user", u, "error", err)
		return
	}

	slog.Debug("owner", "path", filePath, "user", u.Username)
}

func getMostRecentCacheItem() {
	myMap := mycache.Items()
	keys := make([]string, 0, len(myMap))
	for key := range myMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	key := keys[len(keys)-1]

	slog.Debug("cache", "most recent key", key, "value", myMap[key])
}
