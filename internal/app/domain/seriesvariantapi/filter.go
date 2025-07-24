package seriesvariantapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
)

type QueryParams struct {
	SeriesID string
	FuelType string
}

func parseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	filter := QueryParams{
		SeriesID: values.Get("seriesId"),
	}
	return filter
}

func parseQueryFilter(qp QueryParams) (seriesvariant.QueryFilter, error) {
	var qf seriesvariant.QueryFilter

	if qp.SeriesID != "" {
		id, err := strconv.Atoi(qp.SeriesID)
		if err != nil {
			return seriesvariant.QueryFilter{}, fmt.Errorf("invalid brand id: %w", err)
		}
		qf.SeriesID = &id
	}

	return qf, nil
}
