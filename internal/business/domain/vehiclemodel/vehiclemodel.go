package vehiclemodel

import (
	"context"
	"errors"
	"fmt"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var (
	ErrNotfound                = errors.New("vehicle model not found")
	ErrConflict                = errors.New("provided input conflicts with existing data")
	ErrInvalidDriveType        = errors.New("invalid drive type")
	ErrInvalidFuelType         = errors.New("invalid fuel type")
	ErrInvalidBodyStyle        = errors.New("invalid body style")
	ErrInvalidTransmissionType = errors.New("invalid transmission style")
	ErrInvalidEngineType       = errors.New("invalid engine type")
	ErrDuplicatedModel         = errors.New("duplicated model")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]VehicleModel, error)
	QueryByID(ctx context.Context, id int) (VehicleModel, error)
	Create(ctx context.Context, vm VehicleModel) (VehicleModel, error)
	Update(ctx context.Context, vm VehicleModel) (VehicleModel, error)
	Delete(ctx context.Context, id int) error
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) NewWithTx(txM tran.TxManager) (*Core, error) {
	storer, err := c.storer.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &Core{
		storer: storer,
	}, nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]VehicleModel, error) {
	vs, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return vs, nil
}

func (c *Core) QueryByID(ctx context.Context, id int) (VehicleModel, error) {
	vm, err := c.storer.QueryByID(ctx, id)
	if err != nil {
		return VehicleModel{}, fmt.Errorf("querybyid: %w", err)
	}
	return vm, nil
}

func (c *Core) Create(ctx context.Context, nvm NewVehicleModel) (VehicleModel, error) {
	vm := VehicleModel{
		SeriesName:         nvm.SeriesName,
		CommercialName:     nvm.CommercialName,
		ModelYear:          nvm.ModelYear,
		BrandID:            nvm.BrandID,
		Nickname:           nvm.Nickname,
		EngineDisplacement: nvm.EngineDisplacement,
		DriveType:          nvm.DriveType,
		FuelType:           nvm.FuelType,
		BodyStyle:          nvm.BodyStyle,
		TransmissionType:   nvm.TransmissionType,
	}

	vm, err := c.storer.Create(ctx, vm)
	if err != nil {
		return VehicleModel{}, fmt.Errorf("create: %w", err)
	}

	return vm, nil
}

func (c *Core) Update(ctx context.Context, vm VehicleModel, uvm UpdateVehicleModel) (VehicleModel, error) {
	if uvm.SeriesName != nil {
		vm.SeriesName = *uvm.SeriesName
	}

	if uvm.CommercialName != nil {
		vm.CommercialName = *uvm.CommercialName
	}

	if uvm.ModelYear != nil {
		vm.ModelYear = *uvm.ModelYear
	}

	if uvm.Nickname != nil {
		vm.ModelYear = *uvm.ModelYear
	}

	if uvm.EngineDisplacement != nil {
		vm.EngineDisplacement = *uvm.EngineDisplacement
	}

	if uvm.DriveType != nil {
		vm.DriveType = *uvm.DriveType
	}

	if uvm.FuelType != nil {
		vm.FuelType = *uvm.FuelType
	}

	if uvm.BodyStyle != nil {
		vm.BodyStyle = *uvm.BodyStyle
	}

	if uvm.TransmissionType != nil {
		vm.TransmissionType = *uvm.TransmissionType
	}

	vm, err := c.storer.Update(ctx, vm)
	if err != nil {
		return VehicleModel{}, fmt.Errorf("update: %w", err)
	}

	vm, err = c.storer.QueryByID(ctx, vm.ID)
	if err != nil {
		return VehicleModel{}, fmt.Errorf("querybyid: %w", err)
	}

	return vm, nil
}

func (c *Core) Delete(ctx context.Context, id int) error {
	if err := c.storer.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
