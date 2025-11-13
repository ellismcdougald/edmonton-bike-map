package updater

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

const overpassURL = "https://overpass-api.de/api/interpreter"

func FetchOverpassData(query []byte) ([]byte, error) {
	resp, err := http.Post(overpassURL, "application/x-www-form-urlencoded", bytes.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer func() {
    if cerr := resp.Body.Close(); cerr != nil {
        log.Printf("failed to close response body: %v", cerr)
    }
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}