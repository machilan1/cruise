package brandseriesapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/business/domain/brandseries"
)

type QueryParams struct{}

func parseQueryParams(r *http.Request) QueryParams {
	// values := r.URL.Query()

	filter := QueryParams{}
	return filter
}

func parseQueryFilter(qp QueryParams) (brandseries.QueryFilter, error) {
	return brandseries.QueryFilter{}, nil
}
