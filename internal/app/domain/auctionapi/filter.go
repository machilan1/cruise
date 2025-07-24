package auctionapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/business/domain/auction"
)

type QueryParams struct {
	Page     string
	PageSize string
	OrderBy  string
}

func parseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	filter := QueryParams{
		Page:     values.Get("page"),
		PageSize: values.Get("pageSize"),
		OrderBy:  values.Get("orderBy"),
	}

	return filter
}

func parseQueryFilter(qp QueryParams) (auction.QueryFilter, error) {
	deleted := false
	filter := auction.QueryFilter{
		Deleted: &deleted,
	}

	return filter, nil
}
