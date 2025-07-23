package seriesmodelapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
)

type QueryParams struct {
	BrandID string
}

func parseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	filter := QueryParams{
		BrandID: values.Get("brandId"),
	}
	return filter
}

func parseQueryFilter(qp QueryParams) (seriesmodel.QueryFilter, error) {
	var qf seriesmodel.QueryFilter

	if qp.BrandID != "" {
		id, err := strconv.Atoi(qp.BrandID)
		if err != nil {
			return seriesmodel.QueryFilter{}, fmt.Errorf("invalid brand id: %w", err)
		}
		qf.BrandID = &id
	}

	return qf, nil
}
