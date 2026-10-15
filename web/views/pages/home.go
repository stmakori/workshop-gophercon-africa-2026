package pages

import (
	"workshop/internal/domain/sighting"
	"workshop/web/views/components"
)

// entryPosition is a sighting's place in the results, or 0 when it was off route.
// The list is ordered by distance, so every sighting on the route comes first.
func entryPosition(s sighting.Sighting, i int) int {
	if !components.OnRoute(s) {
		return 0
	}
	return i + 1
}
