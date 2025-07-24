package seriesvariantapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	seriesModelKey ctxKey = iota
)

func seriesVariantCtx(smc *seriesvariant.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "seriesVariantID"); id != "" {
				dID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				brd, err := smc.QueryByID(ctx, dID)
				if err != nil {
					if errors.Is(err, seriesvariant.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: seriesVariantID[%d]: %w", dID, err)
				}
				ctx = setSeriesVariant(ctx, brd)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func setSeriesVariant(ctx context.Context, sm seriesvariant.SeriesVariant) context.Context {
	return context.WithValue(ctx, seriesModelKey, sm)
}

func getSeriesVariant(ctx context.Context) (seriesvariant.SeriesVariant, error) {
	dsh, ok := ctx.Value(seriesModelKey).(seriesvariant.SeriesVariant)
	if !ok {
		return seriesvariant.SeriesVariant{}, fmt.Errorf("series variant not found in context")
	}

	return dsh, nil
}
