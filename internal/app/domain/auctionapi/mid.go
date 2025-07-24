package auctionapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	auctionKey ctxKey = iota
)

func getAuction(ctx context.Context) (auction.Auction, error) {
	auc, ok := ctx.Value(auctionKey).(auction.Auction)
	if !ok {
		return auction.Auction{}, fmt.Errorf("auction not found in context")
	}

	return auc, nil
}

func setAuction(ctx context.Context, auc auction.Auction) context.Context {
	return context.WithValue(ctx, auctionKey, auc)
}

func auctionCtx(aucCore *auction.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "auctionID"); id != "" {
				aucID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				auc, err := aucCore.QueryByID(ctx, aucID)
				if err != nil {
					if errors.Is(err, auction.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: auctionID[%d]: %w", aucID, err)
				}
				ctx = setAuction(ctx, auc)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}
