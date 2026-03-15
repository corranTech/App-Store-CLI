package web

import (
	"context"
	"net/http"
)

const (
	integrationsAPIRefererURL           = appStoreBaseURL + "/access/integrations/api"
	integrationsIndividualKeysRefererURL = appStoreBaseURL + "/access/integrations/api/individual-keys"
)

func integrationsHeaders(referer string) http.Header {
	headers := make(http.Header)
	headers.Set("Accept", "application/vnd.api+json, application/json, text/csv")
	headers.Set("Content-Type", "application/json")
	headers.Set("X-CSRF-ITC", "[asc-ui]")
	headers.Set("Origin", appStoreBaseURL)
	headers.Set("Referer", referer)
	return headers
}

func olympusHeaders(referer string) http.Header {
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Requested-With", "xsdr2$")
	if referer != "" {
		headers.Set("Referer", referer)
	}
	return headers
}

func (c *Client) doIrisV1Request(ctx context.Context, method, path string, body any) ([]byte, error) {
	return c.doRequestBase(ctx, irisV1BaseURL, method, path, body, integrationsHeaders(integrationsAPIRefererURL))
}

func (c *Client) doIrisV2Request(ctx context.Context, method, path string, body any) ([]byte, error) {
	return c.doRequestBase(ctx, irisV2BaseURL, method, path, body, integrationsHeaders(integrationsIndividualKeysRefererURL))
}

func (c *Client) doOlympusRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	return c.doRequestBase(ctx, olympusBaseURL, method, path, body, olympusHeaders(integrationsIndividualKeysRefererURL))
}
