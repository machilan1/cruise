package brandapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/business/domain/brand"
)

type QueryParams struct{}

func parseQueryParams(r *http.Request) QueryParams {
	// values := r.URL.Query()

	filter := QueryParams{}
	return filter
}

func parseQueryFilter(qp QueryParams) (brand.QueryFilter, error) {
	return brand.QueryFilter{}, nil
}
