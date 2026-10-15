package components

import (
	"testing"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"

	"github.com/stretchr/testify/assert"
)

func at(lat, lng float64) sighting.Sighting {
	return sighting.Sighting{Location: sighting.Location{Coordinates: sighting.Coordinates{Latitude: lat, Longitude: lng}}}
}

func TestDistanceKm(t *testing.T) {
	t.Parallel()

	// Nairobi to Naivasha town is about 77 km as the crow flies.
	assert.InDelta(t, 77, distanceKm(-1.2864, 36.8172, -0.7167, 36.4333), 3)
	assert.InDelta(t, 0, distanceKm(-0.9, 36.3, -0.9, 36.3), 1e-9)
}

func TestStageLabel(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		s       sighting.Sighting
		want    string
		onRoute bool
		stage   string
	}{
		"On a stage": {
			s:       at(-0.9167, 36.3167),
			want:    "0.0 km from Hell's Gate",
			onRoute: true,
			stage:   "Hell's Gate",
		},
		"Near the lake": {
			s:       at(-0.75, 36.40),
			want:    "5.2 km from Naivasha service park",
			onRoute: true,
			stage:   "Naivasha service park",
		},
		"Off route": {
			s:       at(-1.2864, 36.8172),
			want:    "Off route · 58 km from Kedong",
			onRoute: false,
			stage:   "Off route",
		},
		"No position": {
			s:       sighting.Sighting{},
			want:    "Position not recorded",
			onRoute: false,
			stage:   "Off route",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, StageLabel(test.s))
			assert.Equal(t, test.onRoute, OnRoute(test.s))
			assert.Equal(t, test.stage, StageName(test.s))
		})
	}
}

func TestByDistance(t *testing.T) {
	t.Parallel()

	far, near, none := at(3.5, 39.0), at(-0.92, 36.32), sighting.Sighting{}
	got := ByDistance([]sighting.Sighting{none, far, near})
	assert.Equal(t, []sighting.Sighting{near, far, none}, got)
}

func TestLicence(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "CC BY-NC 4.0", licence("http://creativecommons.org/licenses/by-nc/4.0/"))
	assert.Equal(t, "Public domain", licence("https://creativecommons.org/publicdomain/zero/1.0/"))
	assert.Equal(t, "All rights reserved", licence("All rights reserved"))
}

func TestDisplayName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Hippo", Animal(sighting.Sighting{Species: species.Species{VernacularName: "hippo"}}))
	assert.Equal(t, "Unknown species", Animal(sighting.Sighting{}))
}
