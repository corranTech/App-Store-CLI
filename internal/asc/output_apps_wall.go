package asc

import "strings"

// AppWallEntry represents one row in apps wall output.
type AppWallEntry struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Platform    []string `json:"platform,omitempty"`
	Creator     string   `json:"creator,omitempty"`
	AppStoreURL string   `json:"appStoreUrl"`
	IconURL     string   `json:"iconUrl,omitempty"`
	ReleaseDate string   `json:"releaseDate,omitempty"`
}

// AppsWallResult is the response payload for apps wall output.
type AppsWallResult struct {
	Data []AppWallEntry `json:"data"`
}

func appsWallRows(resp *AppsWallResult) ([]string, [][]string) {
	headers := []string{"App", "Link", "Creator", "Platform", "Icon"}
	rows := make([][]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		rows = append(rows, []string{
			compactWhitespace(item.Name),
			item.AppStoreURL,
			compactWhitespace(item.Creator),
			strings.Join(item.Platform, ", "),
			item.IconURL,
		})
	}
	return headers, rows
}
