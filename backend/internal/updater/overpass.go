package updater

import (
	"bytes"
	"io"
	"net/http"
)

const overpassURL = "https://overpass-api.de/api/interpreter"

func FetchOverpassData(query []byte) ([]byte, error) {
	resp, err := http.Post(overpassURL, "application/x-www-form-urlencoded", bytes.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}