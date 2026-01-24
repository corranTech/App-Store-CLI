package asc

// BundleIDAttributes describes a bundle ID resource.
type BundleIDAttributes struct {
	Name       string   `json:"name"`
	Identifier string   `json:"identifier"`
	Platform   Platform `json:"platform"`
	SeedID     string   `json:"seedId,omitempty"`
}

// BundleIDCreateAttributes describes attributes for creating a bundle ID.
type BundleIDCreateAttributes struct {
	Name       string   `json:"name"`
	Identifier string   `json:"identifier"`
	Platform   Platform `json:"platform"`
}

// BundleIDUpdateAttributes describes attributes for updating a bundle ID.
type BundleIDUpdateAttributes struct {
	Name string `json:"name,omitempty"`
}

// BundleIDCreateData is the data portion of a bundle ID create request.
type BundleIDCreateData struct {
	Type       ResourceType             `json:"type"`
	Attributes BundleIDCreateAttributes `json:"attributes"`
}

// BundleIDCreateRequest is a request to create a bundle ID.
type BundleIDCreateRequest struct {
	Data BundleIDCreateData `json:"data"`
}

// BundleIDUpdateData is the data portion of a bundle ID update request.
type BundleIDUpdateData struct {
	Type       ResourceType              `json:"type"`
	ID         string                    `json:"id"`
	Attributes *BundleIDUpdateAttributes `json:"attributes,omitempty"`
}

// BundleIDUpdateRequest is a request to update a bundle ID.
type BundleIDUpdateRequest struct {
	Data BundleIDUpdateData `json:"data"`
}

// BundleIDCapabilityAttributes describes a bundle ID capability resource.
type BundleIDCapabilityAttributes struct {
	CapabilityType string              `json:"capabilityType"`
	Settings       []CapabilitySetting `json:"settings,omitempty"`
}

// BundleIDCapabilityCreateAttributes describes attributes for creating a capability.
type BundleIDCapabilityCreateAttributes struct {
	CapabilityType string              `json:"capabilityType"`
	Settings       []CapabilitySetting `json:"settings,omitempty"`
}

// CapabilitySetting describes a capability setting.
type CapabilitySetting struct {
	Key     string             `json:"key"`
	Name    string             `json:"name,omitempty"`
	Options []CapabilityOption `json:"options,omitempty"`
}

// CapabilityOption describes a capability option.
type CapabilityOption struct {
	Key         string `json:"key"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
}

// BundleIDCapabilityRelationships describes relationships for bundle ID capabilities.
type BundleIDCapabilityRelationships struct {
	BundleID *Relationship `json:"bundleId"`
}

// BundleIDCapabilityCreateData is the data portion of a capability create request.
type BundleIDCapabilityCreateData struct {
	Type          ResourceType                       `json:"type"`
	Attributes    BundleIDCapabilityCreateAttributes `json:"attributes"`
	Relationships *BundleIDCapabilityRelationships   `json:"relationships"`
}

// BundleIDCapabilityCreateRequest is a request to create a bundle ID capability.
type BundleIDCapabilityCreateRequest struct {
	Data BundleIDCapabilityCreateData `json:"data"`
}

// BundleIDsResponse is the response from bundle IDs list endpoint.
type BundleIDsResponse = Response[BundleIDAttributes]

// BundleIDResponse is the response from bundle ID detail endpoint.
type BundleIDResponse = SingleResponse[BundleIDAttributes]

// BundleIDCapabilitiesResponse is the response from bundle ID capabilities endpoint.
type BundleIDCapabilitiesResponse = Response[BundleIDCapabilityAttributes]

// BundleIDCapabilityResponse is the response from bundle ID capability detail endpoint.
type BundleIDCapabilityResponse = SingleResponse[BundleIDCapabilityAttributes]
