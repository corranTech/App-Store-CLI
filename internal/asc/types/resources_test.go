package types

import "testing"

func TestResponseAccessors(t *testing.T) {
	r := &Response[struct{ Name string }]{
		Data: []Resource[struct{ Name string }]{
			{
				Type: ResourceTypeApps,
				ID:   "app-1",
				Attributes: struct{ Name string }{
					Name: "Example",
				},
			},
		},
		Links: Links{
			Self: "/v1/apps",
			Next: "/v1/apps?page=2",
		},
	}

	links := r.GetLinks()
	if links == nil || links.Next != "/v1/apps?page=2" {
		t.Fatalf("unexpected links: %+v", links)
	}

	data, ok := r.GetData().([]Resource[struct{ Name string }])
	if !ok {
		t.Fatalf("expected []Resource data type, got %T", r.GetData())
	}
	if len(data) != 1 || data[0].ID != "app-1" {
		t.Fatalf("unexpected data payload: %+v", data)
	}
}

func TestLinkagesResponseAccessors(t *testing.T) {
	r := &LinkagesResponse{
		Data: []ResourceData{
			{Type: ResourceTypeBuilds, ID: "build-1"},
		},
		Links: Links{
			Self: "/v1/builds",
		},
	}

	links := r.GetLinks()
	if links == nil || links.Self != "/v1/builds" {
		t.Fatalf("unexpected links: %+v", links)
	}

	data, ok := r.GetData().([]ResourceData)
	if !ok {
		t.Fatalf("expected []ResourceData type, got %T", r.GetData())
	}
	if len(data) != 1 || data[0].ID != "build-1" {
		t.Fatalf("unexpected linkage payload: %+v", data)
	}
}

func TestTypeConstants(t *testing.T) {
	if PlatformIOS != "IOS" || PlatformMacOS != "MAC_OS" {
		t.Fatalf("unexpected platform constants: %q %q", PlatformIOS, PlatformMacOS)
	}
	if ChecksumAlgorithmSHA256 != "SHA_256" {
		t.Fatalf("unexpected checksum algorithm constant: %q", ChecksumAlgorithmSHA256)
	}
	if UTIIPA != "com.apple.ipa" || UTIPKG != "com.apple.installer-package-archive" {
		t.Fatalf("unexpected UTI constants: %q %q", UTIIPA, UTIPKG)
	}
}
