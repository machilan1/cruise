package seriesmodel

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
	ErrNotFound                = errors.New("series model not found")
	ErrConflict                = errors.New("input data conflicts with existing data")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]SeriesModel, error)
	QueryByID(ctx context.Context, smID int) (SeriesModel, error)
	Create(ctx context.Context, sm SeriesModel) (SeriesModel, error)
	Update(ctx context.Context, sm SeriesModel) (SeriesModel, error)
	Delete(ctx context.Context, smID int) error
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

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]SeriesModel, error) {
	sms, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return sms, nil
}

func (c *Core) QueryByID(ctx context.Context, smID int) (SeriesModel, error) {
	sm, err := c.storer.QueryByID(ctx, smID)
	if err != nil {
		return SeriesModel{}, fmt.Errorf("querybyid: %w", err)
	}

	return sm, nil
}

func (c *Core) Create(ctx context.Context, nsm NewSeriesModel) (SeriesModel, error) {
	now := time.Now()

	version := ""
	if nsm.Version != nil {
		version = *nsm.Version
	}

	sm := SeriesModel{
		Name:               nsm.Name,
		Version:            version,
		ModelYear:          nsm.ModelYear,
		BodyStyle:          nsm.BodyStyle,
		DriveType:          nsm.DriveType,
		FuelType:           nsm.FuelType,
		EngineType:         nsm.EngineType,
		EngineDisplacement: nsm.EngineDisplacement,
		ValveCount:         nsm.ValveCount,
		HasTurbo:           nsm.HasTurbo,
		TransmissionType:   nsm.TransmissionType,
		HorsePower:         nsm.HorsePower,
		Series: SeriesModelSeries{
			ID: nsm.SeriesID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	sm, err := c.storer.Create(ctx, sm)
	if err != nil {
		return SeriesModel{}, fmt.Errorf("create: %w", err)
	}

	return sm, nil
}

func (c *Core) Update(ctx context.Context, usm UpdateSeriesModel, sm SeriesModel) (SeriesModel, error) {
	if usm.Name != nil {
		sm.Name = *usm.Name
	}
	if usm.Version != nil {
		sm.Version = *usm.Version
	}
	if usm.ModelYear != nil {
		sm.ModelYear = *usm.ModelYear
	}
	if usm.BodyStyle != nil {
		sm.BodyStyle = *usm.BodyStyle
	}
	if usm.DriveType != nil {
		sm.DriveType = *usm.DriveType
	}
	if usm.FuelType != nil {
		sm.FuelType = *usm.FuelType
	}
	if usm.EngineType != nil {
		sm.EngineType = *usm.EngineType
	}
	if usm.EngineDisplacement != nil {
		sm.EngineDisplacement = *usm.EngineDisplacement
	}
	if usm.ValveCount != nil {
		sm.ValveCount = *usm.ValveCount
	}
	if usm.HasTurbo != nil {
		sm.HasTurbo = *usm.HasTurbo
	}
	if usm.TransmissionType != nil {
		sm.TransmissionType = *usm.TransmissionType
	}
	if usm.HorsePower != nil {
		sm.HorsePower = *usm.HorsePower
	}
	if usm.SeriesID != nil {
		sm.Series.ID = *usm.SeriesID
	}

	sm.UpdatedAt = time.Now()
	sm, err := c.storer.Update(ctx, sm)
	if err != nil {
		return SeriesModel{}, fmt.Errorf("update: %w", err)
	}

	return sm, nil
}

func (c *Core) Delete(ctx context.Context, smID int) error {
	if err := c.storer.Delete(ctx, smID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
