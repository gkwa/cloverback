package cloverback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/atotto/clipboard"
)

func writeBufferToClipboard(buffer bytes.Buffer) {
	clipboard.WriteAll(buffer.String())
}

func writeBufferToStdout(buffer bytes.Buffer) {
	_, copyErr := io.Copy(os.Stdout, &buffer)
	if copyErr != nil {
		slog.Error("copy to stdout", "error", copyErr.Error())
	}
}

func saveValuesToFile() {
	myMap := mycache.Items()
	keys := make([]string, 0, len(myMap))
	for key := range myMap {
		keys = append(keys, key)
	}
	pattern := `^(.+)_([0-9]+)$`
	re := regexp.MustCompile(pattern)

	// tempDir := os.TempDir()
	tempDir := "/tmp/cloverback"
	err := os.MkdirAll(tempDir, 0o755)
	if err != nil {
		slog.Error("make temp dir", "error", err.Error())
		panic(err)
	}
	for _, key := range keys {
		submatches := re.FindStringSubmatch(key)

		if len(submatches) == 3 {
			key1 := submatches[1]
			timestamp := submatches[2]
			slog.Debug("regex", "key", key1, "timestamp", timestamp)

			fname := fmt.Sprintf("%s.json", timestamp)

			tempFile := filepath.Join(tempDir, fname)

			slog.Debug("save file", "fname", tempFile)

			if _, err := os.Stat(tempFile); err == nil {
				slog.Debug("File exists")
			} else if os.IsNotExist(err) {
				file, err := os.Create(tempFile)
				if err != nil {
					slog.Error("file create", "error", err.Error())
					return
				}
				defer file.Close()

				v, _ := mycache.Get(key)

				var x PushbulletHTTReply
				bytes := []byte(v.(string))
				err = json.Unmarshal(bytes, &x)
				if err != nil {
					slog.Error("unmarshalling", "error", err.Error())
					panic(err)
				}

				jsBytes, err := json.MarshalIndent(x, "", "  ")
				if err != nil {
					panic(err)
				}

				_, err = file.WriteString(string(jsBytes))
				if err != nil {
					slog.Error("file write", "error", err.Error())
					return
				}
				slog.Debug("Config file created successfully!")

			} else {
				slog.Error("file check", "error", err.Error())
			}
		} else {
			slog.Debug("Input does not match the expected pattern.")
		}
	}
}
