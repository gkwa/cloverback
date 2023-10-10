package cloverback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/atotto/clipboard"
)

func deleteAllPushbulletRecords() {
	apiKey := os.Getenv("PUSHBULLET_API_KEY")
	if apiKey == "" {
		slog.Error("PUSHBULLET_API_KEY environment variable is not set.")
	}

	apiUrl := "https://api.pushbullet.com/v2/pushes"

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", apiUrl, nil)
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
		slog.Debug("Request was successful. Pushes deleted.")
	} else {
		fmt.Printf("Request failed with status code %d\n", resp.StatusCode)
	}
}

func Main() {
	apiKey := os.Getenv("PUSHBULLET_API_KEY")
	if apiKey == "" {
		slog.Error("PUSHBULLET_API_KEY environment variable is not set.")
	}

	apiURL := "https://api.pushbullet.com/v2/pushes"
	client := &http.Client{}

	queryParams := map[string]string{
		"active":         "true",
		"modified_after": "1.4e+09",
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		slog.Error("Error creating HTTP request", "error", err.Error())
		return
	}

	req.Header.Set("Access-Token", apiKey)

	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending GET request", "error", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("pushbullet response", "status_code", resp.StatusCode)
		return
	}
	slog.Debug("pushbullet response", "status_code", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response", "error", err.Error())
		return
	}

	pushBulletJsonBlob := string(body)
	cacheString(pushBulletJsonBlob)

	var result PushbulletHTTReply
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		slog.Error("unmarshalling response body", "error", err.Error())
		return
	}

	x, err := json.Marshal(result)
	slog.Debug("json response", "json", x)
	if err != nil {
		slog.Error("marhsalling response body", "error", err.Error())
		return
	}

	cachePath, err := getCachePath(getCacheRelPath())
	if err != nil {
		panic(err)
	}

	pbPushesCache.LoadFile(cachePath)

	cacheString(string(x))

	if len(result.Pushes) > 0 {
		OutputOrgMode(os.Stdout, result)
	}

	slog.Debug("caching", "cache path", cachePath)
	pbPushesCache.SaveFile(cachePath)
	slog.Debug("cache", "item count", pbPushesCache.ItemCount())
	getMostRecent()
	saveValuesToFile()
	deleteAllPushbulletRecords()
}

func OutputOrgMode(output io.Writer, reply PushbulletHTTReply) {
	tmplStr := `{{range .Pushes}}
*** {{.Title}}
**** summary 

{{.URL}}
{{end}}`

	tmpl, err := template.New("pushTemplate").Parse(tmplStr)
	if err != nil {
		slog.Error("parsing template", "error", err.Error())
		return
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(output, reply)
	if err != nil {
		slog.Error("executing template", "error", err.Error())
		return
	}

	err = tmpl.Execute(&outputBuffer, reply)
	if err != nil {
		fmt.Println("executing template error:", err)
		return
	}

	clipboard.WriteAll(outputBuffer.String())

	_, copyErr := io.Copy(output, &outputBuffer)
	if copyErr != nil {
		fmt.Println("copying error:", copyErr)
		return
	}
}

func saveValuesToFile() {
	myMap := pbPushesCache.Items()
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

				v, _ := pbPushesCache.Get(key)

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
