package seriesmodelapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	seriesModelKey ctxKey = iota
)

func seriesModelCtx(smc *seriesmodel.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "seriesModelID"); id != "" {
				dID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				brd, err := smc.QueryByID(ctx, dID)
				if err != nil {
					if errors.Is(err, seriesmodel.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: seriesModelID[%d]: %w", dID, err)
				}
				ctx = setSeriesModel(ctx, brd)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func setSeriesModel(ctx context.Context, sm seriesmodel.SeriesModel) context.Context {
	return context.WithValue(ctx, seriesModelKey, sm)
}

func getSeriesModel(ctx context.Context) (seriesmodel.SeriesModel, error) {
	dsh, ok := ctx.Value(seriesModelKey).(seriesmodel.SeriesModel)
	if !ok {
		return seriesmodel.SeriesModel{}, fmt.Errorf("series model not found in context")
	}

	return dsh, nil
}
