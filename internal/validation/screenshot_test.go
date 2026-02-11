package validation

import "testing"

func TestScreenshotChecks_Mismatch(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-1",
			DisplayType: "APP_IPHONE_65",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-1", FileName: "shot.png", Width: 100, Height: 100},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if !hasCheckID(checks, "screenshots.dimension_mismatch") {
		t.Fatalf("expected dimension mismatch check")
	}
}

func TestScreenshotChecks_Pass(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-1",
			DisplayType: "APP_IPHONE_65",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-1", FileName: "shot.png", Width: 1242, Height: 2688},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if len(checks) != 0 {
		t.Fatalf("expected no checks, got %d", len(checks))
	}
}
