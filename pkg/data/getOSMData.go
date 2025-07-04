package data

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

/*
var query string = `
[out:json][timeout:30];

// Find Canada area by ISO3166 code (country)
area["ISO3166-1"="CA"][admin_level=2]->.canada;

// Find Edmonton area inside Canada (admin_level=8)
area["name"="Edmonton"][admin_level=8](area.canada)->.edmonton_area;

// Query ways inside Edmonton area
(
  way["highway"]["area"!~"yes"]["highway"!~"motorway|motorway_link|raceway|construction|service"]["bicycle"!~"no"](area.edmonton_area);
  way["highway"="cycleway"]["bicycle"!~"no"](area.edmonton_area);
);
out body;
>;
out skel qt;
`
*/

func GetOSMData(query string) error {
	resp, err := http.Post("https://overpass-api.de/api/interpreter", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte("data="+query)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile("osm_bike_data.json", body, 0644)
	if err != nil {
		return err
	}
	return nil
}