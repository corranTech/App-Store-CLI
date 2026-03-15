package builds

import (
	"testing"
)

func TestDsymsCommandShape(t *testing.T) {
	cmd := BuildsDsymsCommand()
	if cmd == nil {
		t.Fatal("expected dsyms command")
	}
	if cmd.Name != "dsyms" {
		t.Errorf("expected name dsyms, got %s", cmd.Name)
	}

	flagNames := []string{"build", "output-dir", "output"}
	for _, name := range flagNames {
		if cmd.FlagSet.Lookup(name) == nil {
			t.Errorf("expected flag --%s to be registered", name)
		}
	}
}

func TestDsymsRequiresBuild(t *testing.T) {
	cmd := BuildsDsymsCommand()
	err := cmd.Exec(t.Context(), nil)
	if err == nil {
		t.Fatal("expected error for missing --build")
	}
}

func TestExtractDSYMURLs(t *testing.T) {
	s := func(v string) *string { return &v }

	tests := []struct {
		name     string
		bundles  []dsymBundleInfo
		wantURLs int
	}{
		{
			name:     "no bundles",
			bundles:  nil,
			wantURLs: 0,
		},
		{
			name: "bundle with dsym url",
			bundles: []dsymBundleInfo{
				{BundleID: "com.example.app", DSYMURL: s("https://example.com/dsym.zip")},
			},
			wantURLs: 1,
		},
		{
			name: "bundle without dsym url",
			bundles: []dsymBundleInfo{
				{BundleID: "com.example.app", DSYMURL: nil},
			},
			wantURLs: 0,
		},
		{
			name: "mixed bundles",
			bundles: []dsymBundleInfo{
				{BundleID: "com.example.app", DSYMURL: s("https://example.com/app.zip")},
				{BundleID: "com.example.clip", DSYMURL: nil},
				{BundleID: "com.example.ext", DSYMURL: s("https://example.com/ext.zip")},
			},
			wantURLs: 2,
		},
		{
			name: "empty dsym url string",
			bundles: []dsymBundleInfo{
				{BundleID: "com.example.app", DSYMURL: s("")},
			},
			wantURLs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterBundlesWithDSYM(tt.bundles)
			if len(got) != tt.wantURLs {
				t.Errorf("got %d bundles with dsym URLs, want %d", len(got), tt.wantURLs)
			}
		})
	}
}

func TestDsymFileName(t *testing.T) {
	tests := []struct {
		bundleID string
		buildID  string
		index    int
		want     string
	}{
		{"com.example.app", "build-1", 0, "com.example.app.dsym.zip"},
		{"", "build-1", 0, "build-1_0.dsym.zip"},
		{"", "build-1", 2, "build-1_2.dsym.zip"},
	}

	for _, tt := range tests {
		got := dsymFileName(tt.bundleID, tt.buildID, tt.index)
		if got != tt.want {
			t.Errorf("dsymFileName(%q, %q, %d) = %q, want %q", tt.bundleID, tt.buildID, tt.index, got, tt.want)
		}
	}
}
