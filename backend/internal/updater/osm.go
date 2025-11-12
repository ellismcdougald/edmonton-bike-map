package updater

import "encoding/json"

// OSMResponse represents the structure of the response from the OpenStreetMap API.
type OSMResponse struct {
	Elements []OSMElement `json:"elements"`
}

// OSMElement represents a single element (node or way) in the OSM data.
type OSMElement struct {
	Type  string             `json:"type"`
	ID    int64              `json:"id"`
	Lat   float64            `json:"lat,omitempty"`
	Lon   float64            `json:"lon,omitempty"`
	Tags  map[string]string  `json:"tags,omitempty"`
	Nodes []int64            `json:"nodes,omitempty"`
}

// ParseOSMBytes unmarshals OSM data directly from bytes.
func ParseOSMBytes(data []byte) (*OSMResponse, error) {
	var resp OSMResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}