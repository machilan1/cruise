package vehiclemodelapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidID       = errs.NewTrustedError(fmt.Errorf("invalid vehicle model id"), http.StatusBadRequest)
	ErrNotFound        = errs.NewTrustedError(fmt.Errorf("vehicle model not found"), http.StatusNotFound)
	ErrConflict        = errs.NewTrustedError(fmt.Errorf("request data conflict with current data"), http.StatusConflict)
	ErrDuplicatedModel = errs.NewTrustedError(fmt.Errorf("duplicate model"), http.StatusConflict)
)

type handlers struct {
	log          *logger.Logger
	txM          tran.TxManager
	vehicleModel *vehiclemodel.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, vehicleModelCore *vehiclemodel.Core) *handlers {
	return &handlers{
		log:          log,
		txM:          txM,
		vehicleModel: vehicleModelCore,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.vehicleModel.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:          h.log,
		txM:          txM,
		vehicleModel: fl,
	}, nil
}

func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := parseQueryParams(r)
	filter, err := parseQueryFilter(qp)
	if err != nil {
		return err
	}

	vms, err := h.vehicleModel.Query(ctx, filter)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return web.Respond(ctx, w, toAppVehicleModels(vms), http.StatusOK)
}

func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vm, err := getVehicleModel(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppVehicleModel(vm), http.StatusOK)
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewVehicleModel
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	var vm vehiclemodel.VehicleModel
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		nvm, err := toCoreNewVehicleModel(app)
		if err != nil {
			return err
		}

		vm, err = h.vehicleModel.Create(ctx, nvm)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		vm, err = h.vehicleModel.QueryByID(ctx, vm.ID)
		if err != nil {
			return fmt.Errorf("querybyid: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, vehiclemodel.ErrConflict) {
			return ErrConflict
		}

		if errors.Is(err, vehiclemodel.ErrDuplicatedModel) {
			return ErrDuplicatedModel
		}

		return err
	}
	return web.Respond(ctx, w, toAppVehicleModel(vm), http.StatusCreated)
}

func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vm, err := getVehicleModel(ctx)
	if err != nil {
		return err
	}

	var app AppUpdateVehicleModel
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		uvm, err := toCoreUpdateVehicleModel(app)
		if err != nil {
			return err
		}

		vm, err = h.vehicleModel.Update(ctx, vm, uvm)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, vehiclemodel.ErrDuplicatedModel) {
			return ErrDuplicatedModel
		}

		if errors.Is(err, vehiclemodel.ErrConflict) {
			return ErrConflict
		}

		if errors.Is(err, vehiclemodel.ErrNotfound) {
			return ErrNotFound
		}
		return err
	}

	return web.Respond(ctx, w, toAppVehicleModel(vm), http.StatusOK)
}

func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vm, err := getVehicleModel(ctx)
	if err != nil {
		return err
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		if err := h.vehicleModel.Delete(ctx, vm.ID); err != nil {
			return fmt.Errorf("delete: %w", err)
		}
		return nil
	}); err != nil {
		if errors.Is(err, vehiclemodel.ErrConflict) {
			return ErrConflict
		}

		if errors.Is(err, vehiclemodel.ErrNotfound) {
			return ErrNotFound
		}

		return err
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
