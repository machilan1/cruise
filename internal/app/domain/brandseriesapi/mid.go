package brandseriesapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/brandseries"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	brandSeriesKey ctxKey = iota
)

func brandSeriesCtx(bs brandseries.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "brandSeriesID"); id != "" {
				dID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				brd, err := bs.QueryByID(ctx, dID)
				if err != nil {
					if errors.Is(err, brandseries.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: brandSeriesID[%d]: %w", dID, err)
				}
				ctx = setBrandSeries(ctx, brd)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func setBrandSeries(ctx context.Context, bs brandseries.BrandSeries) context.Context {
	return context.WithValue(ctx, brandSeriesKey, bs)
}
func getBrandSeries(ctx context.Context) (brandseries.BrandSeries, error) {
	dsh, ok := ctx.Value(brandSeriesKey).(brandseries.BrandSeries)
	if !ok {
		return brandseries.BrandSeries{}, fmt.Errorf("brand series not found in context")
	}

	return dsh, nil
}
