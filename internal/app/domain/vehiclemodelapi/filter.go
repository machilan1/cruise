package vehiclemodelapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
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

func parseQueryFilter(qp QueryParams) (vehiclemodel.QueryFilter, error) {
	var filter vehiclemodel.QueryFilter

	if qp.BrandID != "" {
		id, err := strconv.Atoi(qp.BrandID)
		if err != nil {
			return vehiclemodel.QueryFilter{}, fmt.Errorf("invalid brand id")
		}
		filter.BrandID = &id
	}

	return filter, nil
}
