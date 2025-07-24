package auctionapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/app/sdk/query"
	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/sdk/order"
	"github.com/machilan1/cruise/internal/business/sdk/paging"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidID           = errs.NewTrustedError(fmt.Errorf("invalid auction id"), http.StatusBadRequest)
	ErrNotFound            = errs.NewTrustedError(fmt.Errorf("auction not found"), http.StatusNotFound)
	ErrEtagVersionConflict = errs.NewTrustedError(fmt.Errorf("etag version conflict"), http.StatusPreconditionFailed)
	ErrConflict            = errs.NewTrustedError(fmt.Errorf("request data conflict with current data"), http.StatusConflict)
)

type handlers struct {
	log     *logger.Logger
	txM     tran.TxManager
	auction *auction.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, auction *auction.Core) *handlers {
	return &handlers{
		log:     log,
		txM:     txM,
		auction: auction,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	auc, err := h.auction.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:     h.log,
		txM:     txM,
		auction: auc,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)

	page, err := paging.Parse(qp.Page, qp.PageSize)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, auction.DefaultOrderBy)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	qf, err := parseQueryFilter(qp)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	aucs, err := h.auction.Query(ctx, qf, orderBy, page)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.auction.Count(ctx, qf)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, query.NewResult(toAppAuctions(aucs), total, page), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	auc, err := getAuction(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppAuction(auc), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewAuction
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	nAuc, err := toCoreNewAuction(app)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var auc auction.Auction
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		auc, err = h.auction.Create(ctx, nAuc)
		if err != nil {
			return fmt.Errorf("create: auc[%+v]: %w", app, err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, auction.ErrDataConflict) {
			return ErrConflict
		}
		return err
	}

	return web.Respond(ctx, w, toAppAuction(auc), http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateAuction
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	uAuc, err := toCoreUpdateAuction(app)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	auc, err := getAuction(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		auc, err = h.auction.Update(ctx, auc, uAuc)
		if err != nil {
			return fmt.Errorf("update: auctionID[%d] app[%+v]: %w", auc.ID, app, err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, auction.ErrDataConflict) {
			return ErrConflict
		}
		return err
	}

	return web.Respond(ctx, w, toAppAuction(auc), http.StatusOK)
}

func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	auc, err := getAuction(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.auction.Delete(ctx, auc); err != nil {
			return fmt.Errorf("delete: auctionID[%d]: %w", auc.ID, err)
		}

		return nil
	}); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) softDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// TODO: If restore is ever needed, please isolate a new queryByID without using view to filter deletedAt, otherwise the precondition is not available.
	auc, err := getAuction(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		auc, err = h.auction.Archive(ctx, auc)
		if err != nil {
			return fmt.Errorf("delete: auctionID[%d]: %w", auc.ID, err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, auction.ErrDataConflict) {
			return ErrConflict
		}
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
