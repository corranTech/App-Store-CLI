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

// CertificateAttributes describes a certificate resource.
type CertificateAttributes struct {
	Name               string `json:"name"`
	CertificateType    string `json:"certificateType"`
	DisplayName        string `json:"displayName,omitempty"`
	SerialNumber       string `json:"serialNumber,omitempty"`
	Platform           string `json:"platform,omitempty"`
	ExpirationDate     string `json:"expirationDate,omitempty"`
	CertificateContent string `json:"certificateContent,omitempty"`
}

// CertificateCreateAttributes describes attributes for creating a certificate.
type CertificateCreateAttributes struct {
	CertificateType string `json:"certificateType"`
	CSRContent      string `json:"csrContent"`
}

// CertificateCreateData is the data portion of a certificate create request.
type CertificateCreateData struct {
	Type       ResourceType                `json:"type"`
	Attributes CertificateCreateAttributes `json:"attributes"`
}

// CertificateCreateRequest is a request to create a certificate.
type CertificateCreateRequest struct {
	Data CertificateCreateData `json:"data"`
}

// CertificatesResponse is the response from certificates list endpoint.
type CertificatesResponse = Response[CertificateAttributes]

// CertificateResponse is the response from certificate detail endpoint.
type CertificateResponse = SingleResponse[CertificateAttributes]

// DeviceAttributes describes a device resource.
type DeviceAttributes struct {
	Name        string   `json:"name"`
	UDID        string   `json:"udid"`
	Platform    Platform `json:"platform"`
	Status      string   `json:"status"`
	DeviceClass string   `json:"deviceClass,omitempty"`
	Model       string   `json:"model,omitempty"`
	AddedDate   string   `json:"addedDate,omitempty"`
}

// DeviceCreateAttributes describes attributes for creating a device.
type DeviceCreateAttributes struct {
	Name     string   `json:"name"`
	UDID     string   `json:"udid"`
	Platform Platform `json:"platform"`
}

// DeviceUpdateAttributes describes attributes for updating a device.
type DeviceUpdateAttributes struct {
	Status string `json:"status,omitempty"`
}

// DeviceCreateData is the data portion of a device create request.
type DeviceCreateData struct {
	Type       ResourceType           `json:"type"`
	Attributes DeviceCreateAttributes `json:"attributes"`
}

// DeviceCreateRequest is a request to create a device.
type DeviceCreateRequest struct {
	Data DeviceCreateData `json:"data"`
}

// DeviceUpdateData is the data portion of a device update request.
type DeviceUpdateData struct {
	Type       ResourceType            `json:"type"`
	ID         string                  `json:"id"`
	Attributes *DeviceUpdateAttributes `json:"attributes,omitempty"`
}

// DeviceUpdateRequest is a request to update a device.
type DeviceUpdateRequest struct {
	Data DeviceUpdateData `json:"data"`
}

// DevicesResponse is the response from devices list endpoint.
type DevicesResponse = Response[DeviceAttributes]

// DeviceResponse is the response from device detail endpoint.
type DeviceResponse = SingleResponse[DeviceAttributes]
