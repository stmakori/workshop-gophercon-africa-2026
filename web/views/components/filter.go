package components

import (
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"workshop/internal/domain/sighting"
)

// Animal names the animal seen, as the card does.
func Animal(s sighting.Sighting) string { return displayName(s.Species) }

// Year is when, to the year.
func Year(s sighting.Sighting) string { return strconv.Itoa(s.HappenedAt.Year()) }

// FilterOptions counts how often each key appears, keeping the n most common; 0 keeps all.
func FilterOptions(sightings []sighting.Sighting, n int, key func(sighting.Sighting) string) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, key), n)
}

// YearOptions lists the years with sightings, newest first, for the year chips.
func YearOptions(sightings []sighting.Sighting) []sighting.Tally {
	rows := FilterOptions(sightings, 0, Year)
	slices.SortFunc(rows, func(a, b sighting.Tally) int { return strings.Compare(b.Name, a.Name) })
	return rows
}

// StageOptions lists the stages with sightings, most first, then "Off route".
func StageOptions(sightings []sighting.Sighting) []sighting.Tally {
	rows := FilterOptions(sightings, 0, StageName)
	i := slices.IndexFunc(rows, func(t sighting.Tally) bool { return t.Name == "Off route" })
	if i >= 0 {
		off := rows[i]
		rows = append(slices.Delete(rows, i, i+1), off)
	}
	return rows
}

// FilterAttrs marks a list item with what the filter matches on.
func FilterAttrs(s sighting.Sighting) templ.Attributes {
	return templ.Attributes{
		"data-animal": Animal(s),
		"data-year":   Year(s),
		"data-stage":  StageName(s),
		"data-search": strings.ToLower(strings.Join([]string{
			Animal(s), s.Species.ScientificName, place(s.Location), s.RecordedBy,
		}, " ")),
	}
}
