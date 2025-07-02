package data

import (
    "encoding/csv"
    "os"
		"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
		"fmt"
		"strconv"
    "strings"
)

func getGeometryLine(coordinateStr string) []model.Coordinate {
	coordinateStr = strings.TrimPrefix(coordinateStr, "MULTILINESTRING ")
	coordinateStr = strings.Trim(coordinateStr, "()")
	coordinatePairs := strings.Split(coordinateStr, ",")

	var allCoordinates []model.Coordinate
	for _, pair := range coordinatePairs {
		pair = strings.TrimSpace(pair)
		coords := strings.Split(pair, " ")
		latitude, _ := strconv.ParseFloat(coords[0], 64)
		longitude, _ := strconv.ParseFloat(coords[1], 64)

		allCoordinates = append(allCoordinates, model.Coordinate{
			Latitude: latitude,
			Longitude: longitude,
		})
	}

	return allCoordinates
}

func GetRouteData(fileName string) []model.Route {
	var rawData [][]string = readCsv(fileName)

	var routeData []model.Route
	for i, values := range rawData {
		//var allCoords []model.Coordinate
		fmt.Println(values[9])
		getGeometryLine(values[9])
		/*
		for j, raw_coords := range values[9] {
			coords := model.Coordinate {
				Latitude: raw_coords[0],
				Longitude: raw_coords[1],
			}
			allCoords[j] = coords
		}

		route := model.Route {
			Id: values[0],
			Duration: values[4],
			GeometryLine: allCoords,
		}
		routeData[i] = route
		*/
		if i == 1 { break }
	}

	return routeData
}

func readCsv(fileName string) [][]string {
	f, _ := os.Open(fileName)
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()

	return records
}