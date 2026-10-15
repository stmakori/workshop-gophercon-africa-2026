package handlers_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"
	"workshop/internal/common/logs"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/web/handlers"

	"github.com/stretchr/testify/assert"
)

func TestSighting(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	found := sighting.Sighting{ID: id, Species: species.Species{VernacularName: "plains zebra"}}

	tt := map[string]struct {
		path   string
		finder fakeFinder
		want   int
		body   string
	}{
		"Found": {
			path:   "/sightings/" + id.String(),
			finder: fakeFinder{item: found},
			want:   http.StatusOK,
			body:   "Plains zebra",
		},
		"Not a UUID": {
			path: "/sightings/not-an-id",
			want: http.StatusNotFound,
			body: "wrong stage",
		},
		"Not found": {
			path:   "/sightings/" + uuid.New().String(),
			finder: fakeFinder{err: sighting.ErrNotFound},
			want:   http.StatusNotFound,
			body:   "wrong stage",
		},
		"Store error": {
			path:   "/sightings/" + id.String(),
			finder: fakeFinder{err: errors.New("store unavailable")},
			want:   http.StatusInternalServerError,
			body:   "Unable to load this sighting.",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.Handle("GET /sightings/{id}", handlers.Sighting(discard(), test.finder))
			got := serve(t, mux, test.path)

			assert.Equal(t, test.want, got.Code)
			assert.Equal(t, "text/html; charset=utf-8", got.Header().Get("Content-Type"))
			assert.Contains(t, got.Body.String(), test.body)
		})
	}
}

func TestRoutePages(t *testing.T) {
	t.Parallel()

	pages := map[string]func(*slog.Logger, handlers.SightingsLister) http.HandlerFunc{
		"Stats":   handlers.Stats,
		"Gallery": handlers.Gallery,
	}

	for name, page := range pages {
		t.Run(name+"/OK", func(t *testing.T) {
			t.Parallel()

			lister := fakeLister{items: []sighting.Sighting{
				{ID: uuid.New(), Species: species.Species{VernacularName: "impala"}},
			}}
			got := serve(t, page(discard(), lister), "/")
			assert.Equal(t, http.StatusOK, got.Code)
		})

		t.Run(name+"/No sightings", func(t *testing.T) {
			t.Parallel()

			got := serve(t, page(discard(), fakeLister{}), "/")
			assert.Equal(t, http.StatusOK, got.Code)
			assert.Contains(t, got.Body.String(), "go run . db seed")
		})

		t.Run(name+"/List error", func(t *testing.T) {
			t.Parallel()

			got := serve(t, page(discard(), fakeLister{err: errors.New("store unavailable")}), "/")
			assert.Equal(t, http.StatusInternalServerError, got.Code)
		})
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	got := serve(t, handlers.NotFound(discard()), "/nowhere")
	assert.Equal(t, http.StatusNotFound, got.Code)
	assert.Contains(t, got.Body.String(), "Off the road")
}

type fakeFinder struct {
	item sighting.Sighting
	err  error
}

func (f fakeFinder) Find(context.Context, uuid.UUID) (sighting.Sighting, error) {
	return f.item, f.err
}

func discard() *slog.Logger {
	return logs.NewProduction(io.Discard, slog.LevelError)
}

func serve(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))

	return rec
}
