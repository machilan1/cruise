package brandseriesapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/brandseries"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrConflict  = errs.NewTrustedError(fmt.Errorf("input data conflicts with existing data"), http.StatusConflict)
	ErrNotFound  = errs.NewTrustedError(fmt.Errorf("brand series not found"), http.StatusNotFound)
	ErrInvalidID = errs.NewTrustedError(fmt.Errorf("invalid brand series id"), http.StatusBadRequest)
)

type handlers struct {
	log         *logger.Logger
	txM         tran.TxManager
	brandSeries *brandseries.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, brandSeries *brandseries.Core) *handlers {
	return &handlers{
		log:         log,
		txM:         txM,
		brandSeries: brandSeries,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.brandSeries.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:         h.log,
		txM:         txM,
		brandSeries: fl,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)

	qf, err := parseQueryFilter(qp)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	bs, err := h.brandSeries.Query(ctx, qf)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return web.Respond(ctx, w, toAppBrandSerieses(bs), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	bs, err := getBrandSeries(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppBrandSeries(bs), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var anbs AppNewBrandSeries
	if err := web.Decode(r, &anbs); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var bs brandseries.BrandSeries
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		nbs, err := toCoreNewBrandSeries(anbs)
		if err != nil {
			return err
		}

		bs, err = h.brandSeries.Create(ctx, nbs)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, brandseries.ErrConflict) {
			return ErrConflict
		}
		return err
	}
	return web.Respond(ctx, w, toAppBrandSeries(bs), http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	bs, err := getBrandSeries(ctx)
	if err != nil {
		return err
	}

	var aubs AppUpdateBrandSeries
	if err := web.Decode(r, &aubs); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		ubs, err := toCoreUpdateBrandSeries(aubs)
		if err != nil {
			return err
		}

		bs, err = h.brandSeries.Update(ctx, ubs, bs)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, brandseries.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, brandseries.ErrConflict) {
			return ErrConflict
		}

		return err
	}

	return web.Respond(ctx, w, toAppBrandSeries(bs), http.StatusOK)
}

func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	bs, err := getBrandSeries(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.brandSeries.Delete(ctx, bs.ID); err != nil {
			return fmt.Errorf("delete: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, brandseries.ErrNotFound) {
			return ErrNotFound
		}

		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
