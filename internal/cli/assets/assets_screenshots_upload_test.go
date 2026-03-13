package assets

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestUploadScreenshotsSkipExistingStartsUploadTimeoutAfterChecksumFiltering(t *testing.T) {
	t.Setenv("ASC_TIMEOUT", "200ms")

	filePath := writeAssetsTestPNG(t, t.TempDir(), "01-home.png")
	fileSizeBytes := fileSize(t, filePath)

	origChecksumFunc := screenshotFileChecksumFunc
	screenshotFileChecksumFunc = func(path string) (string, error) {
		time.Sleep(250 * time.Millisecond)
		return computeFileChecksum(path)
	}
	t.Cleanup(func() {
		screenshotFileChecksumFunc = origChecksumFunc
	})

	origTransport := http.DefaultTransport
	http.DefaultTransport = assetsUploadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if err := req.Context().Err(); err != nil {
			return nil, err
		}

		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/v1/appStoreVersionLocalizations/LOC_123/appScreenshotSets":
			return assetsJSONResponse(http.StatusOK, `{"data":[{"type":"appScreenshotSets","id":"set-1","attributes":{"screenshotDisplayType":"APP_IPHONE_65"}}],"links":{}}`)
		case req.Method == http.MethodGet && req.URL.Path == "/v1/appScreenshotSets/set-1/appScreenshots":
			return assetsJSONResponse(http.StatusOK, `{"data":[],"links":{}}`)
		case req.Method == http.MethodGet && req.URL.Path == "/v1/appScreenshotSets/set-1/relationships/appScreenshots":
			return assetsJSONResponse(http.StatusOK, `{"data":[],"links":{}}`)
		case req.Method == http.MethodPost && req.URL.Path == "/v1/appScreenshots":
			body := fmt.Sprintf(`{"data":{"type":"appScreenshots","id":"new-1","attributes":{"uploadOperations":[{"method":"PUT","url":"https://upload.example/new-1","length":%d,"offset":0}]}}}`, fileSizeBytes)
			return assetsJSONResponse(http.StatusCreated, body)
		case req.Method == http.MethodPut && req.URL.Host == "upload.example":
			return assetsJSONResponse(http.StatusOK, `{}`)
		case req.Method == http.MethodPatch && req.URL.Path == "/v1/appScreenshots/new-1":
			return assetsJSONResponse(http.StatusOK, `{"data":{"type":"appScreenshots","id":"new-1","attributes":{"uploaded":true}}}`)
		case req.Method == http.MethodGet && req.URL.Path == "/v1/appScreenshots/new-1":
			return assetsJSONResponse(http.StatusOK, `{"data":{"type":"appScreenshots","id":"new-1","attributes":{"assetDeliveryState":{"state":"COMPLETE"}}}}`)
		case req.Method == http.MethodPatch && req.URL.Path == "/v1/appScreenshotSets/set-1/relationships/appScreenshots":
			return assetsJSONResponse(http.StatusNoContent, "")
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})
	t.Cleanup(func() {
		http.DefaultTransport = origTransport
	})

	client := newAssetsUploadTestClient(t)
	result, err := uploadScreenshots(context.Background(), client, "LOC_123", "APP_IPHONE_65", []string{filePath}, true, false)
	if err != nil {
		t.Fatalf("uploadScreenshots() error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 upload result, got %d", len(result.Results))
	}
	if result.Results[0].AssetID != "new-1" {
		t.Fatalf("expected uploaded asset ID new-1, got %#v", result.Results[0])
	}
}
