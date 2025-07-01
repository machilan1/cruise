package brandapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/brand"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	brandKey ctxKey = iota
)

func brandCtx(b *brand.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "brandID"); id != "" {
				dID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				brd, err := b.QueryByID(ctx, dID)
				if err != nil {
					if errors.Is(err, brand.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: brandID[%d]: %w", dID, err)
				}
				ctx = setBrand(ctx, brd)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func setBrand(ctx context.Context, dsh brand.Brand) context.Context {
	return context.WithValue(ctx, brandKey, dsh)
}

func getBrand(ctx context.Context) (brand.Brand, error) {
	dsh, ok := ctx.Value(brandKey).(brand.Brand)
	if !ok {
		return brand.Brand{}, fmt.Errorf("brand not found in context")
	}

	return dsh, nil
}
