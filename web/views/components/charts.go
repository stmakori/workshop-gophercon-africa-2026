package components

import (
	"maps"
	"slices"
	"strconv"

	"workshop/internal/domain/sighting"
)

// peak finds the busiest year and its count.
func peak(byYear map[int]int) (year, count int) {
	for y, n := range byYear {
		if n > count || (n == count && y < year) {
			year, count = y, n
		}
	}
	return year, count
}

// yearsNewestFirst lists the years that have sightings, latest at the top like a results sheet.
func yearsNewestFirst(byYear map[int]int) []int {
	ys := slices.Sorted(maps.Keys(byYear))
	slices.Reverse(ys)
	return ys
}

// bar sets --bar to a row's share of the longest, for the CSS to draw.
func bar(count, longest int) map[string]string {
	if longest <= 0 {
		return map[string]string{"--bar": "0%"}
	}
	return map[string]string{"--bar": strconv.Itoa(count*100/longest) + "%"}
}

// TopAnimals ranks the animals seen most often.
func TopAnimals(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, Animal), n)
}
