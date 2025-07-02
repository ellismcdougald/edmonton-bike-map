package data

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

/*
query := `
	[out:json][timeout:30];
	area["name"="Edmonton"]["admin_level"="8"]->.a;
	(
		way["highway"]["area"!~"yes"]["highway"!~"motorway|motorway_link|raceway|construction|service"]["bicycle"!~"no"](area.a);
		way["highway"="cycleway"]["bicycle"!~"no"](area.a);
	);
	out body;
	>;
	out skel qt;
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