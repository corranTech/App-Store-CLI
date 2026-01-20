//go:build integration

package asc

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestIntegrationEndpoints(t *testing.T) {
	keyID := os.Getenv("ASC_KEY_ID")
	issuerID := os.Getenv("ASC_ISSUER_ID")
	keyPath := os.Getenv("ASC_PRIVATE_KEY_PATH")
	appID := os.Getenv("ASC_APP_ID")

	if keyID == "" || issuerID == "" || keyPath == "" || appID == "" {
		t.Skip("integration tests require ASC_KEY_ID, ASC_ISSUER_ID, ASC_PRIVATE_KEY_PATH, ASC_APP_ID")
	}

	client, err := NewClient(keyID, issuerID, keyPath)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Run("feedback", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		feedback, err := client.GetFeedback(ctx, appID, WithFeedbackLimit(1))
		if err != nil {
			t.Fatalf("failed to fetch feedback: %v", err)
		}
		if feedback == nil {
			t.Fatal("expected feedback response")
		}
		assertLimit(t, len(feedback.Data), 1)
		assertASCLink(t, feedback.Links.Self)
		assertASCLink(t, feedback.Links.Next)
		if len(feedback.Data) > 0 && feedback.Data[0].Type == "" {
			t.Fatal("expected feedback data type to be set")
		}
		if feedback.Links.Next != "" {
			nextCtx, nextCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer nextCancel()
			nextFeedback, err := client.GetFeedback(nextCtx, appID, WithFeedbackNextURL(feedback.Links.Next))
			if err != nil {
				t.Fatalf("failed to fetch feedback next page: %v", err)
			}
			if nextFeedback == nil {
				t.Fatal("expected feedback next page response")
			}
			assertASCLink(t, nextFeedback.Links.Self)
			assertASCLink(t, nextFeedback.Links.Next)
		}

		if len(feedback.Data) == 0 {
			t.Skip("no feedback data to validate filters")
		}

		first := feedback.Data[0].Attributes
		if first.DeviceModel != "" {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetFeedback(
				filteredCtx,
				appID,
				WithFeedbackDeviceModels([]string{first.DeviceModel}),
				WithFeedbackLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered feedback by device model: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no feedback results for device model filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.DeviceModel != first.DeviceModel {
					t.Fatalf("expected device model %q, got %q", first.DeviceModel, item.Attributes.DeviceModel)
				}
			}
		}

		if first.OSVersion != "" {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetFeedback(
				filteredCtx,
				appID,
				WithFeedbackOSVersions([]string{first.OSVersion}),
				WithFeedbackLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered feedback by os version: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no feedback results for os version filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.OSVersion != first.OSVersion {
					t.Fatalf("expected os version %q, got %q", first.OSVersion, item.Attributes.OSVersion)
				}
			}
		}
	})

	t.Run("crashes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		crashes, err := client.GetCrashes(ctx, appID, WithCrashLimit(1))
		if err != nil {
			t.Fatalf("failed to fetch crashes: %v", err)
		}
		if crashes == nil {
			t.Fatal("expected crashes response")
		}
		assertLimit(t, len(crashes.Data), 1)
		assertASCLink(t, crashes.Links.Self)
		assertASCLink(t, crashes.Links.Next)
		if len(crashes.Data) > 0 && crashes.Data[0].Type == "" {
			t.Fatal("expected crash data type to be set")
		}
		if crashes.Links.Next != "" {
			nextCtx, nextCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer nextCancel()
			nextCrashes, err := client.GetCrashes(nextCtx, appID, WithCrashNextURL(crashes.Links.Next))
			if err != nil {
				t.Fatalf("failed to fetch crashes next page: %v", err)
			}
			if nextCrashes == nil {
				t.Fatal("expected crashes next page response")
			}
			assertASCLink(t, nextCrashes.Links.Self)
			assertASCLink(t, nextCrashes.Links.Next)
		}

		if len(crashes.Data) == 0 {
			t.Skip("no crash data to validate filters")
		}

		first := crashes.Data[0].Attributes
		if first.DeviceModel != "" {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetCrashes(
				filteredCtx,
				appID,
				WithCrashDeviceModels([]string{first.DeviceModel}),
				WithCrashLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered crashes by device model: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no crash results for device model filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.DeviceModel != first.DeviceModel {
					t.Fatalf("expected device model %q, got %q", first.DeviceModel, item.Attributes.DeviceModel)
				}
			}
		}

		if first.OSVersion != "" {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetCrashes(
				filteredCtx,
				appID,
				WithCrashOSVersions([]string{first.OSVersion}),
				WithCrashLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered crashes by os version: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no crash results for os version filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.OSVersion != first.OSVersion {
					t.Fatalf("expected os version %q, got %q", first.OSVersion, item.Attributes.OSVersion)
				}
			}
		}
	})

	t.Run("reviews", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		reviews, err := client.GetReviews(ctx, appID, WithLimit(1))
		if err != nil {
			t.Fatalf("failed to fetch reviews: %v", err)
		}
		if reviews == nil {
			t.Fatal("expected reviews response")
		}
		assertLimit(t, len(reviews.Data), 1)
		assertASCLink(t, reviews.Links.Self)
		assertASCLink(t, reviews.Links.Next)
		if len(reviews.Data) > 0 && reviews.Data[0].Type == "" {
			t.Fatal("expected review data type to be set")
		}
		if reviews.Links.Next != "" {
			nextCtx, nextCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer nextCancel()
			nextReviews, err := client.GetReviews(nextCtx, appID, WithNextURL(reviews.Links.Next))
			if err != nil {
				t.Fatalf("failed to fetch reviews next page: %v", err)
			}
			if nextReviews == nil {
				t.Fatal("expected reviews next page response")
			}
			assertASCLink(t, nextReviews.Links.Self)
			assertASCLink(t, nextReviews.Links.Next)
		}

		if len(reviews.Data) == 0 {
			t.Skip("no review data to validate filters")
		}

		first := reviews.Data[0].Attributes
		if first.Rating >= 1 && first.Rating <= 5 {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetReviews(
				filteredCtx,
				appID,
				WithRating(first.Rating),
				WithLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered reviews by rating: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no review results for rating filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.Rating != first.Rating {
					t.Fatalf("expected rating %d, got %d", first.Rating, item.Attributes.Rating)
				}
			}
		}

		if first.Territory != "" {
			filteredCtx, filteredCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer filteredCancel()
			filtered, err := client.GetReviews(
				filteredCtx,
				appID,
				WithTerritory(first.Territory),
				WithLimit(5),
			)
			if err != nil {
				t.Fatalf("failed to fetch filtered reviews by territory: %v", err)
			}
			assertLimit(t, len(filtered.Data), 5)
			if len(filtered.Data) == 0 {
				t.Skip("no review results for territory filter")
			}
			for _, item := range filtered.Data {
				if item.Attributes.Territory != first.Territory {
					t.Fatalf("expected territory %q, got %q", first.Territory, item.Attributes.Territory)
				}
			}
		}
	})
}

func assertLimit(t *testing.T, count, limit int) {
	t.Helper()
	if limit <= 0 {
		return
	}
	if count > limit {
		t.Fatalf("expected at most %d items, got %d", limit, count)
	}
}

func assertASCLink(t *testing.T, link string) {
	t.Helper()
	if link == "" {
		return
	}
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("expected link to be a valid URL, got %q: %v", link, err)
	}
	if parsed.Host != "" && parsed.Host != "api.appstoreconnect.apple.com" {
		t.Fatalf("expected App Store Connect host, got %q", parsed.Host)
	}
	if parsed.Scheme != "" && parsed.Scheme != "https" {
		t.Fatalf("expected https scheme, got %q", parsed.Scheme)
	}
}
