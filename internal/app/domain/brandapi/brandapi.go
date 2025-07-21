package brandapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/brand"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidID = errs.NewTrustedError(fmt.Errorf("invalid brand id"), http.StatusBadRequest)
	ErrNotFound  = errs.NewTrustedError(fmt.Errorf("brand not found"), http.StatusNotFound)
	ErrConflict  = errs.NewTrustedError(fmt.Errorf("request data conflict with current data"), http.StatusConflict)
)

type handlers struct {
	log   *logger.Logger
	txM   tran.TxManager
	brand *brand.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, brand *brand.Core) *handlers {
	return &handlers{
		log:   log,
		txM:   txM,
		brand: brand,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.brand.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:   h.log,
		txM:   txM,
		brand: fl,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)

	qf, err := parseQueryFilter(qp)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	brds, err := h.brand.Query(ctx, qf)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return web.Respond(ctx, w, toAppBrands(brds), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	b, err := getBrand(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppBrand(b), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var anbrd AppNewBrand
	if err := web.Decode(r, &anbrd); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var brd brand.Brand
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		brd, err = h.brand.Create(ctx, toCoreNewBrand(anbrd))
		if err != nil {
			return fmt.Errorf("create: anbrd[%+v]: %w", anbrd, err)
		}

		brd, err = h.brand.QueryByID(ctx, brd.ID)
		if err != nil {
			return fmt.Errorf("querybyid: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, brand.ErrConflict) {
			return ErrConflict
		}

		if errors.Is(err, brand.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return web.Respond(ctx, w, toAppBrand(brd), http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	brd, err := getBrand(ctx)
	if err != nil {
		return err
	}

	var aubrd AppUpdateBrand
	if err := web.Decode(r, &aubrd); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		brd, err = h.brand.Update(ctx, brd, toCoreUpdateBrand(aubrd))
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		brd, err = h.brand.QueryByID(ctx, brd.ID)
		if err != nil {
			return fmt.Errorf("querybyid: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, brand.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, brand.ErrConflict) {
			return ErrConflict
		}

		return err
	}
	return web.Respond(ctx, w, toAppBrand(brd), http.StatusOK)
}

func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	brd, err := getBrand(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err = h.brand.Delete(ctx, brd.ID); err != nil {
			return fmt.Errorf("delete: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, brand.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, brand.ErrConflict) {
			return ErrConflict
		}

		return err
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
