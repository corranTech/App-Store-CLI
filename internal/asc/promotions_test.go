package asc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestCreateAppStoreVersionPromotion(t *testing.T) {
	resp := AppStoreVersionPromotionResponse{
		Data: Resource[AppStoreVersionPromotionAttributes]{
			Type: ResourceTypeAppStoreVersionPromotions,
			ID:   "promo-123",
		},
	}
	body, _ := json.Marshal(resp)

	client := newTestClient(t, func(req *http.Request) {
		assertAuthorized(t, req)
		if req.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if !strings.HasSuffix(req.URL.Path, "/v1/appStoreVersionPromotions") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}

		var createReq AppStoreVersionPromotionCreateRequest
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if createReq.Data.Type != ResourceTypeAppStoreVersionPromotions {
			t.Errorf("expected type appStoreVersionPromotions, got %s", createReq.Data.Type)
		}
		if createReq.Data.Relationships.AppStoreVersion == nil {
			t.Fatal("expected appStoreVersion relationship to be set")
		}
		if createReq.Data.Relationships.AppStoreVersion.Data.ID != "version-123" {
			t.Errorf("expected version ID version-123, got %s", createReq.Data.Relationships.AppStoreVersion.Data.ID)
		}
		if createReq.Data.Relationships.AppStoreVersionExperimentTreatment != nil {
			t.Errorf("expected treatment relationship to be nil")
		}
	}, jsonResponse(http.StatusCreated, string(body)))

	result, err := client.CreateAppStoreVersionPromotion(context.Background(), "version-123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Data.ID != "promo-123" {
		t.Errorf("expected ID promo-123, got %s", result.Data.ID)
	}
}

func TestCreateAppStoreVersionPromotion_WithTreatment(t *testing.T) {
	resp := AppStoreVersionPromotionResponse{
		Data: Resource[AppStoreVersionPromotionAttributes]{
			Type: ResourceTypeAppStoreVersionPromotions,
			ID:   "promo-456",
		},
	}
	body, _ := json.Marshal(resp)

	client := newTestClient(t, func(req *http.Request) {
		assertAuthorized(t, req)

		var createReq AppStoreVersionPromotionCreateRequest
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if createReq.Data.Relationships.AppStoreVersionExperimentTreatment == nil {
			t.Fatal("expected treatment relationship to be set")
		}
		rel := createReq.Data.Relationships.AppStoreVersionExperimentTreatment
		if rel.Data.ID != "treatment-456" {
			t.Errorf("expected treatment ID treatment-456, got %s", rel.Data.ID)
		}
		if rel.Data.Type != ResourceTypeAppStoreVersionExperimentTreatments {
			t.Errorf("expected type appStoreVersionExperimentTreatments, got %s", rel.Data.Type)
		}
	}, jsonResponse(http.StatusCreated, string(body)))

	result, err := client.CreateAppStoreVersionPromotion(context.Background(), "version-123", "treatment-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Data.ID != "promo-456" {
		t.Errorf("expected ID promo-456, got %s", result.Data.ID)
	}
}
