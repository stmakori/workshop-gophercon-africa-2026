package components

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"workshop/internal/domain/sighting"
)

// Stage is one area the 2026 Safari Rally Kenya ran through. The positions are
// the centre of each area, not the roads the stages took.
type Stage struct {
	Name  string  `json:"name"`
	Day   string  `json:"day"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Power bool    `json:"power,omitzero"`
	Base  bool    `json:"base,omitzero"`
}

// Stages are in running order, from the service park to the Power Stage.
var Stages = []Stage{
	{Name: "Naivasha service park", Day: "Base", Lat: -0.7167, Lng: 36.4333, Base: true},
	{Name: "Loldia", Day: "Friday", Lat: -0.7000, Lng: 36.3200},
	{Name: "KenGen Geothermal", Day: "Friday", Lat: -0.8833, Lng: 36.2917},
	{Name: "Kedong", Day: "Friday", Lat: -0.9200, Lng: 36.4400},
	{Name: "Soysambu", Day: "Saturday", Lat: -0.4700, Lng: 36.1900},
	{Name: "Elmenteita", Day: "Saturday", Lat: -0.4500, Lng: 36.2500},
	{Name: "Sleeping Warrior", Day: "Saturday", Lat: -0.4900, Lng: 36.2400},
	{Name: "Oserengoni", Day: "Sunday", Lat: -0.8000, Lng: 36.2900},
	{Name: "Hell's Gate", Day: "Sunday", Lat: -0.9167, Lng: 36.3167, Power: true},
}

// RouteKm is how close to a stage a sighting must be to count as on the route.
const RouteKm = 25.0

// NearestStage finds the stage closest to a sighting and how far away it is in km.
// It reports false when the sighting has no position.
func NearestStage(s sighting.Sighting) (Stage, float64, bool) {
	if !hasPosition(s) {
		return Stage{}, 0, false
	}
	best, bestKm := Stage{}, math.Inf(1)
	for _, st := range Stages {
		if km := distanceKm(s.Location.Latitude, s.Location.Longitude, st.Lat, st.Lng); km < bestKm {
			best, bestKm = st, km
		}
	}
	return best, bestKm, true
}

// OnRoute is whether a sighting was made within RouteKm of a stage.
func OnRoute(s sighting.Sighting) bool {
	_, km, ok := NearestStage(s)
	return ok && km <= RouteKm
}

// StageLabel says where a sighting sits against the route, such as "4.2 km from Hell's Gate".
func StageLabel(s sighting.Sighting) string {
	st, km, ok := NearestStage(s)
	switch {
	case !ok:
		return "Position not recorded"
	case km <= RouteKm:
		return fmt.Sprintf("%s km from %s", kilometres(km), st.Name)
	default:
		return fmt.Sprintf("Off route · %s km from %s", kilometres(km), st.Name)
	}
}

// StageName is the stage a sighting counts towards, or "Off route".
func StageName(s sighting.Sighting) string {
	st, km, ok := NearestStage(s)
	if !ok || km > RouteKm {
		return "Off route"
	}
	return st.Name
}

// ByDistance orders sightings closest to a stage first, with those that have no position last.
func ByDistance(sightings []sighting.Sighting) []sighting.Sighting {
	km := func(s sighting.Sighting) float64 {
		if _, d, ok := NearestStage(s); ok {
			return d
		}
		return math.Inf(1)
	}
	return slices.SortedStableFunc(slices.Values(sightings), func(a, b sighting.Sighting) int {
		return cmp.Compare(km(a), km(b))
	})
}

// CountOnRoute is how many sightings were made on the route.
func CountOnRoute(sightings []sighting.Sighting) int {
	n := 0
	for _, s := range sightings {
		if OnRoute(s) {
			n++
		}
	}
	return n
}

// StageTally ranks the stages by how many sightings were made on them.
func StageTally(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, func(s sighting.Sighting) string {
		if !OnRoute(s) {
			return ""
		}
		return StageName(s)
	}), n)
}

func hasPosition(s sighting.Sighting) bool {
	return s.Location.Latitude != 0 || s.Location.Longitude != 0
}

// distanceKm is the great-circle distance between two points, by the haversine formula.
func distanceKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthKm = 6371.0
	rad := math.Pi / 180
	dLat, dLng := (lat2-lat1)*rad, (lng2-lng1)*rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * earthKm * math.Asin(math.Sqrt(a))
}

// kilometres renders a distance to one decimal place under 10 km, whole numbers above.
func kilometres(km float64) string {
	if km < 10 {
		return fmt.Sprintf("%.1f", km)
	}
	return fmt.Sprintf("%.0f", km)
}
