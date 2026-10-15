package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/internal/domain/species"
	"workshop/web/views/pages"

	"github.com/a-h/templ"
)

type SightingsLister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
}

// routeFilter picks Kenya's sightings of the animals the Safari Rally passes on its
// Naivasha, Hell's Gate and Soysambu stages, rhinos included for Soysambu.
func routeFilter() sighting.ListFilter {
	return sighting.ListFilter{
		Country: country.Kenya,
		Animals: []species.Animal{
			species.AnimalZebra,
			species.AnimalGiraffe,
			species.AnimalAfricanBuffalo,
			species.AnimalImpala,
			species.AnimalCommonWarthog,
			species.AnimalWildebeest,
			species.AnimalHippopotamus,
			species.AnimalGreaterFlamingo,
			species.AnimalSpottedHyena,
			species.AnimalRhinoceros,
		},
	}
}

func Home(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return routePage(logger, sightings, func(items []sighting.Sighting) templ.Component {
		return pages.Home(items)
	})
}

// routePage lists the route's sightings and renders them with page.
func routePage(logger *slog.Logger, sightings SightingsLister, page func([]sighting.Sighting) templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, routeFilter())
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		render(ctx, w, logger, page(items))
	}
}

// render writes component as an HTML page, falling back to the error page if it fails.
func render(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(ctx, w); err != nil {
		logger.ErrorContext(ctx, "handlers: rendering page: "+err.Error())
		renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
	}
}
