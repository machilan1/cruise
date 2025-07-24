package seriesvariant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var (
	ErrInvalidDriveType        = errors.New("invalid drive type")
	ErrInvalidFuelType         = errors.New("invalid fuel type")
	ErrInvalidBodyStyle        = errors.New("invalid body style")
	ErrInvalidTransmissionType = errors.New("invalid transmission style")
	ErrInvalidEngineType       = errors.New("invalid engine type")
	ErrNotFound                = errors.New("series variant not found")
	ErrConflict                = errors.New("input data conflicts with existing data")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]SeriesVariant, error)
	QueryByID(ctx context.Context, svID int) (SeriesVariant, error)
	Create(ctx context.Context, sv SeriesVariant) (SeriesVariant, error)
	Update(ctx context.Context, sv SeriesVariant) (SeriesVariant, error)
	Delete(ctx context.Context, svID int) error
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

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]SeriesVariant, error) {
	smv, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return smv, nil
}

func (c *Core) QueryByID(ctx context.Context, svID int) (SeriesVariant, error) {
	sv, err := c.storer.QueryByID(ctx, svID)
	if err != nil {
		return SeriesVariant{}, fmt.Errorf("querybyid: %w", err)
	}

	return sv, nil
}

func (c *Core) Create(ctx context.Context, nsv NewSeriesVariant) (SeriesVariant, error) {
	now := time.Now()

	version := ""
	if nsv.Version != nil {
		version = *nsv.Version
	}

	sm := SeriesVariant{
		Name:               nsv.Name,
		Version:            version,
		ModelYear:          nsv.ModelYear,
		BodyStyle:          nsv.BodyStyle,
		DriveType:          nsv.DriveType,
		FuelType:           nsv.FuelType,
		EngineType:         nsv.EngineType,
		EngineDisplacement: nsv.EngineDisplacement,
		ValveCount:         nsv.ValveCount,
		HasTurbo:           nsv.HasTurbo,
		TransmissionType:   nsv.TransmissionType,
		HorsePower:         nsv.HorsePower,
		Series: SeriesVariantSeries{
			ID: nsv.SeriesID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	sm, err := c.storer.Create(ctx, sm)
	if err != nil {
		return SeriesVariant{}, fmt.Errorf("create: %w", err)
	}

	return sm, nil
}

func (c *Core) Update(ctx context.Context, usv UpdateSeriesVariant, sv SeriesVariant) (SeriesVariant, error) {
	if usv.Name != nil {
		sv.Name = *usv.Name
	}
	if usv.Version != nil {
		sv.Version = *usv.Version
	}
	if usv.ModelYear != nil {
		sv.ModelYear = *usv.ModelYear
	}
	if usv.BodyStyle != nil {
		sv.BodyStyle = *usv.BodyStyle
	}
	if usv.DriveType != nil {
		sv.DriveType = *usv.DriveType
	}
	if usv.FuelType != nil {
		sv.FuelType = *usv.FuelType
	}
	if usv.EngineType != nil {
		sv.EngineType = *usv.EngineType
	}
	if usv.EngineDisplacement != nil {
		sv.EngineDisplacement = *usv.EngineDisplacement
	}
	if usv.ValveCount != nil {
		sv.ValveCount = *usv.ValveCount
	}
	if usv.HasTurbo != nil {
		sv.HasTurbo = *usv.HasTurbo
	}
	if usv.TransmissionType != nil {
		sv.TransmissionType = *usv.TransmissionType
	}
	if usv.HorsePower != nil {
		sv.HorsePower = *usv.HorsePower
	}
	if usv.SeriesID != nil {
		sv.Series.ID = *usv.SeriesID
	}

	sv.UpdatedAt = time.Now()
	sv, err := c.storer.Update(ctx, sv)
	if err != nil {
		return SeriesVariant{}, fmt.Errorf("update: %w", err)
	}

	return sv, nil
}

func (c *Core) Delete(ctx context.Context, svID int) error {
	if err := c.storer.Delete(ctx, svID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
