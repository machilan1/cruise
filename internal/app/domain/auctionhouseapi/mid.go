package auctionhouseapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	auctionHouseKey ctxKey = iota
)

func auctionHouseCtx(ahc *auctionhouse.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "auctionHouseID"); id != "" {
				mID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				ah, err := ahc.QueryByID(ctx, mID)
				if err != nil {
					if errors.Is(err, auctionhouse.ErrNotFound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: auctionHouseID[%d]: %w", mID, err)
				}
				ctx = setAuctionHouse(ctx, ah)
			}

			return next(ctx, w, r)
		}
		return h
	}

	return m
}

func setAuctionHouse(ctx context.Context, ah auctionhouse.AuctionHouse) context.Context {
	return context.WithValue(ctx, auctionHouseKey, ah)
}

func getAuctionHouse(ctx context.Context) (auctionhouse.AuctionHouse, error) {
	vm, ok := ctx.Value(auctionHouseKey).(auctionhouse.AuctionHouse)
	if !ok {
		return auctionhouse.AuctionHouse{}, fmt.Errorf("auction house not found in context")
	}

	return vm, nil
}
