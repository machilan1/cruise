package auctionhouseapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidID = errs.NewTrustedError(fmt.Errorf("invalid auction house id"), http.StatusBadRequest)
	ErrNotFound  = errs.NewTrustedError(fmt.Errorf("auction house not found"), http.StatusNotFound)
	ErrConflict  = errs.NewTrustedError(fmt.Errorf("input data conflicts with existing data"), http.StatusConflict)
)

type handlers struct {
	log          *logger.Logger
	txM          tran.TxManager
	auctionHouse *auctionhouse.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, auctionHouse *auctionhouse.Core) *handlers {
	return &handlers{
		log:          log,
		txM:          txM,
		auctionHouse: auctionHouse,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.auctionHouse.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:          h.log,
		txM:          txM,
		auctionHouse: fl,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)
	filter, err := parseQueryFilter(qp)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	ahs, err := h.auctionHouse.Query(ctx, filter)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return web.Respond(ctx, w, toAppAuctionHouses(ahs), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ah, err := getAuctionHouse(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppAuctionHouse(ah), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nah AppNewAuctionHouse
	if err := web.Decode(r, &nah); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var ah auctionhouse.AuctionHouse
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		na, err := toCoreNewAuctionHouse(nah)
		if err != nil {
			return err
		}

		ah, err = h.auctionHouse.Create(ctx, na)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, auctionhouse.ErrConflict) {
			return ErrConflict
		}
		return err
	}

	return web.Respond(ctx, w, ah, http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ah, err := getAuctionHouse(ctx)
	if err != nil {
		return err
	}

	var uah AppUpdateAuctionHouse
	if err := web.Decode(r, &uah); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		cuah, err := toCoreUpdateAuctionHouse(uah)
		if err != nil {
			return err
		}

		ah, err = h.auctionHouse.Update(ctx, cuah, ah)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		if errors.Is(err, auctionhouse.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, auctionhouse.ErrConflict) {
			return ErrConflict
		}

		return err
	}

	return web.Respond(ctx, w, toAppAuctionHouse(ah), http.StatusOK)
}

func (h *handlers) Archive(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ah, err := getAuctionHouse(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.auctionHouse.Archive(ctx, ah.ID); err != nil {
			return err
		}

		return nil

	}); err != nil {
		if errors.Is(err, auctionhouse.ErrNotFound) {
			return ErrNotFound
		}

		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
