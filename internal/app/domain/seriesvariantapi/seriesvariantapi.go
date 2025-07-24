package seriesvariantapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidBodyStyle        = errs.NewTrustedError(fmt.Errorf("invalid body style"), http.StatusBadRequest)
	ErrInvalidDriveType        = errs.NewTrustedError(fmt.Errorf("invalid drive type"), http.StatusBadRequest)
	ErrInvalidFuelType         = errs.NewTrustedError(fmt.Errorf("invalid fuel type"), http.StatusBadRequest)
	ErrInvalidEngineType       = errs.NewTrustedError(fmt.Errorf("invalid engine type"), http.StatusBadRequest)
	ErrInvalidTransmissionType = errs.NewTrustedError(fmt.Errorf("invalid transmission type"), http.StatusBadRequest)
	ErrConflict                = errs.NewTrustedError(fmt.Errorf("input data conflicts with existing data"), http.StatusConflict)
	ErrNotFound                = errs.NewTrustedError(fmt.Errorf("series variant not found"), http.StatusNotFound)
	ErrInvalidID               = errs.NewTrustedError(fmt.Errorf("invalid series variant id"), http.StatusBadRequest)
)

type handlers struct {
	log           *logger.Logger
	txM           tran.TxManager
	seriesVariant *seriesvariant.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, seriesModel *seriesvariant.Core) *handlers {

	return &handlers{
		log:           log,
		txM:           txM,
		seriesVariant: seriesModel,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.seriesVariant.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:           h.log,
		txM:           txM,
		seriesVariant: fl,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)

	filter, err := parseQueryFilter(qp)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	sms, err := h.seriesVariant.Query(ctx, filter)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return web.Respond(ctx, w, toAppSeriesModels(sms), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sm, err := getSeriesVariant(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppSeriesModel(sm), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ansm AppNewSeriesModel
	if err := web.Decode(r, &ansm); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	nsm, err := toCoreNewSeriesModel(ansm)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var sm seriesvariant.SeriesVariant
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		sm, err = h.seriesVariant.Create(ctx, nsm)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		return nil

	}); err != nil {
		if errors.Is(err, seriesvariant.ErrConflict) {
			return ErrConflict
		}
		return err
	}

	return web.Respond(ctx, w, toAppSeriesModel(sm), http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sm, err := getSeriesVariant(ctx)
	if err != nil {
		return err
	}

	var ausm AppUpdateSeriesModel
	if err := web.Decode(r, &ausm); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	usm, err := toCoreUpdateSeriesModel(ausm)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		sm, err = h.seriesVariant.Update(ctx, usm, sm)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, seriesvariant.ErrConflict) {
			return ErrConflict
		}

		if errors.Is(err, seriesvariant.ErrNotFound) {
			return ErrNotFound
		}

		return err
	}

	return web.Respond(ctx, w, toAppSeriesModel(sm), http.StatusOK)
}

func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sm, err := getSeriesVariant(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.seriesVariant.Delete(ctx, sm.ID); err != nil {
			return fmt.Errorf("delete: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, seriesvariant.ErrNotFound) {
			return ErrNotFound
		}

		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
