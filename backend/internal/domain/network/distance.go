package network

import "math"

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const earthRadius = 6371000 // meters
    degToRad := func(deg float64) float64 { return deg * math.Pi / 180 }

    dLat := degToRad(lat2 - lat1)
    dLon := degToRad(lon2 - lon1)
    lat1Rad := degToRad(lat1)
    lat2Rad := degToRad(lat2)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return earthRadius * c
}
