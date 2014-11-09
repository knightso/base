package utils

import (
	"math"
)

const (
	earthRadius      float64 = 6371 // approximation in kilometers assuming a spherical earth.
	degreesToRadians float64 = math.Pi / 180.0
)

// a unit of return value is kilometer.
func Haversine(lat1, lng1, lat2, lng2 float64) float64 {
	phi1 := lat1 * (degreesToRadians)
	phi2 := lat2 * (degreesToRadians)
	deltaPhi := (lat2 - lat1) * (degreesToRadians)
	deltaLambda := (lng2 - lng1) * (degreesToRadians)

	v1 := math.Sin(deltaPhi / 2)
	v2 := math.Sin(deltaLambda / 2)
	a := v1*v1 + math.Cos(phi1)*math.Cos(phi2)*v2*v2
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
