package components

import (
	"cmp"
	"strings"

	"workshop/internal/domain/sighting"
)

// TopPlaces ranks where along the route the sightings were made.
func TopPlaces(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, func(s sighting.Sighting) string {
		return cmp.Or(s.Location.Locality, s.Location.StateProvince)
	}), n)
}

// TopCrews ranks the observers who recorded the most sightings.
func TopCrews(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, func(s sighting.Sighting) string {
		return s.RecordedBy
	}), n)
}

// PhotoCount is how many photos the wall will show.
func PhotoCount(sightings []sighting.Sighting) int {
	n := 0
	for _, s := range sightings {
		n += len(photos(s.Media))
	}
	return n
}

// photographer credits whoever took a photo, or says nobody was named.
func photographer(m sighting.Media) string {
	return "Photo: " + cmp.Or(m.Creator, m.RightsHolder, "photographer not named")
}

// licence shortens a Creative Commons licence URL to its usual name, such as "CC BY-NC 4.0".
func licence(l string) string {
	_, rest, ok := strings.Cut(l, "creativecommons.org/")
	if !ok {
		return l
	}
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) >= 3 && parts[0] == "licenses" {
		return "CC " + strings.ToUpper(parts[1]) + " " + parts[2]
	}
	if len(parts) >= 2 && parts[0] == "publicdomain" {
		return "Public domain"
	}
	return l
}

// wallCaption is what the lightbox shows under a photo.
func wallCaption(s sighting.Sighting, m sighting.Media) string {
	return displayName(s.Species) + " · " + photographer(m)
}

// photoAlt describes a photo for screen readers.
func photoAlt(s sighting.Sighting, m sighting.Media) string {
	return cmp.Or(m.Title, displayName(s.Species)+", "+place(s.Location))
}

// Place is where a sighting was made, as the card shows it.
func Place(s sighting.Sighting) string { return place(s.Location) }

// FirstPhoto is the photo a sighting leads with.
func FirstPhoto(s sighting.Sighting) (sighting.Media, bool) { return firstPhoto(s.Media) }

// MorePhotos is up to n photos after the first.
func MorePhotos(s sighting.Sighting, n int) []sighting.Media { return morePhotos(s.Media, n) }

// Photographer credits whoever took a photo.
func Photographer(m sighting.Media) string { return photographer(m) }
