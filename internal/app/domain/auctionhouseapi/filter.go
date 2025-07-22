package auctionhouseapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
)

type QueryParams struct {
	CityID string
}

func parseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	filter := QueryParams{
		CityID: values.Get("cityId"),
	}

	return filter
}

func parseQueryFilter(qp QueryParams) (auctionhouse.QueryFilter, error) {
	var filter auctionhouse.QueryFilter

	if qp.CityID != "" {
		id, err := strconv.Atoi(qp.CityID)
		if err != nil {
			return auctionhouse.QueryFilter{}, fmt.Errorf("invalid city id")
		}
		filter.CityID = &id
	}

	return filter, nil

}
