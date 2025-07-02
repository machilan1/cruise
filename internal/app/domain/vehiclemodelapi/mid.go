package vehiclemodelapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
	"github.com/machilan1/cruise/internal/framework/web"
)

type ctxKey int

const (
	vehicleModelKey ctxKey = iota
)

func vehicleModelCtx(vmc *vehiclemodel.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if id := web.Param(r, "vehicleModelID"); id != "" {
				mID, err := strconv.Atoi(id)
				if err != nil {
					return ErrInvalidID
				}

				vehicleModel, err := vmc.QueryByID(ctx, mID)
				if err != nil {
					if errors.Is(err, vehiclemodel.ErrNotfound) {
						return ErrNotFound
					}
					return fmt.Errorf("query: vehicleModelID[%d]: %w", mID, err)
				}
				ctx = setVehicleModel(ctx, vehicleModel)
			}

			return next(ctx, w, r)
		}
		return h
	}

	return m
}

func setVehicleModel(ctx context.Context, vm vehiclemodel.VehicleModel) context.Context {
	return context.WithValue(ctx, vehicleModelKey, vm)
}

func getVehicleModel(ctx context.Context) (vehiclemodel.VehicleModel, error) {
	vm, ok := ctx.Value(vehicleModelKey).(vehiclemodel.VehicleModel)
	if !ok {
		return vehiclemodel.VehicleModel{}, fmt.Errorf("vehicle model not found in context")
	}

	return vm, nil
}
