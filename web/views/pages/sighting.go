package pages

import (
	"fmt"
	"strings"

	"workshop/internal/domain/sighting"
	"workshop/web/views/components"
)

func animalName(s sighting.Sighting) string { return components.Animal(s) }

func where(s sighting.Sighting) string { return components.Place(s) }

func firstPhoto(s sighting.Sighting) (sighting.Media, bool) { return components.FirstPhoto(s) }

// morePhotos is every photo after the first, up to a dozen.
func morePhotos(s sighting.Sighting) []sighting.Media { return components.MorePhotos(s, 12) }

func photoCredit(m sighting.Media) string { return components.Photographer(m) }

// countText says how many were seen, or that GBIF didn't say.
func countText(s sighting.Sighting) string {
	switch {
	case s.IndividualCount == nil:
		return "–"
	case *s.IndividualCount == 1:
		return "1"
	default:
		return fmt.Sprint(*s.IndividualCount)
	}
}

// basis turns GBIF's HUMAN_OBSERVATION into "Human observation".
func basis(b string) string {
	b = strings.ToLower(strings.ReplaceAll(b, "_", " "))
	return strings.ToUpper(b[:1]) + b[1:]
}

// position renders coordinates the way a road book does.
func position(s sighting.Sighting) string {
	lat, lng := s.Location.Latitude, s.Location.Longitude
	ns, ew := "N", "E"
	if lat < 0 {
		ns, lat = "S", -lat
	}
	if lng < 0 {
		ew, lng = "W", -lng
	}
	return fmt.Sprintf("%.4f° %s, %.4f° %s", lat, ns, lng, ew)
}
