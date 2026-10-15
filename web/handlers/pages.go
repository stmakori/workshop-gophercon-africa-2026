package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"uuid"
	"workshop/internal/domain/sighting"
	"workshop/web/views/pages"

	"github.com/a-h/templ"
)

type SightingFinder interface {
	Find(ctx context.Context, id uuid.UUID) (sighting.Sighting, error)
}

// Sighting renders the stage report for one sighting.
func Sighting(logger *slog.Logger, sightings SightingFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			renderNotFound(ctx, w, logger)
			return
		}

		s, err := sightings.Find(ctx, id)
		if errors.Is(err, sighting.ErrNotFound) {
			renderNotFound(ctx, w, logger)
			return
		}
		if err != nil {
			logger.ErrorContext(ctx, "handlers: finding sighting: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load this sighting.")
			return
		}

		render(ctx, w, logger, pages.Sighting(s))
	}
}

// Stats renders the championship standings for the route's sightings.
func Stats(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return routePage(logger, sightings, func(items []sighting.Sighting) templ.Component {
		return pages.Stats(items)
	})
}

// Gallery renders a photo wall of the route's sightings.
func Gallery(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return routePage(logger, sightings, func(items []sighting.Sighting) templ.Component {
		return pages.Gallery(items)
	})
}

// NotFound answers any URL no other route matches.
func NotFound(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderNotFound(r.Context(), w, logger)
	}
}

func renderNotFound(ctx context.Context, w http.ResponseWriter, logger *slog.Logger) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	if err := pages.NotFound().Render(ctx, w); err != nil {
		logger.ErrorContext(ctx, "handlers: rendering not found page: "+err.Error())
	}
}
