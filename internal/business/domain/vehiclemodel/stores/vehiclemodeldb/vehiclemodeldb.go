package vehiclemodeldb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

type Store struct {
	db *sqldb.DB
}

func NewStore(db *sqldb.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) NewWithTx(txM tran.TxManager) (vehiclemodel.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter vehiclemodel.QueryFilter) ([]vehiclemodel.VehicleModel, error) {
	data := map[string]any{
		"brand_id": filter.BrandID,
	}

	const q = `
		SELECT
			vehicle_model_id,
			model_series_name,
			model_commercial_name,
			model_year,
			brand_id,
			nickname,
			engine_displacement,
			drive_type,
			fuel_type,
			body_style,
			transmission_type
		FROM vehicle_models
	`

	var sb strings.Builder

	sb.WriteString(q)
	applyFilter(filter, data, &sb)

	var dbvms []dbVehicleModel
	if err := sqldb.NamedQuerySlice(ctx, s.db, sb.String(), data, &dbvms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreVehicleModels(dbvms)
}

func (s *Store) QueryByID(ctx context.Context, id int) (vehiclemodel.VehicleModel, error) {
	data := map[string]any{
		"vehicle_model_id": id,
	}

	const q = `
		SELECT
			vehicle_model_id,
			model_series_name,
			model_commercial_name,
			model_year,
			brand_id,
			nickname,
			engine_displacement,
			drive_type,
			fuel_type,
			body_style,
			transmission_type
		FROM vehicle_models
		WHERE vehicle_model_id = :vehicle_model_id
	`

	var dbvm dbVehicleModel
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbvm); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return vehiclemodel.VehicleModel{}, vehiclemodel.ErrNotfound
		}
		return vehiclemodel.VehicleModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreVehicleModel(dbvm)
}

func (s *Store) Create(ctx context.Context, vm vehiclemodel.VehicleModel) (vehiclemodel.VehicleModel, error) {
	dbvm := toDBVehicleModel(vm)

	const q = `
		INSERT INTO vehicle_models (
			model_series_name,
			model_commercial_name,
			model_year,
			brand_id,
			nickname,
			engine_displacement,
			drive_type,
			fuel_type,
			body_style,
			transmission_type
		)VALUES(
			:model_series_name,
			:model_commercial_name,
			:model_year,
			:brand_id,
			:nickname,
			:engine_displacement,
			:drive_type,
			:fuel_type,
			:body_style,
			:transmission_type
		)
		RETURNING vehicle_model_id
	`

	var pgErr *pgconn.PgError
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbvm, &dbvm); err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "uq_brand_id_series_name_model_year_model_commercial_name" {
				return vehiclemodel.VehicleModel{}, vehiclemodel.ErrDuplicatedModel
			}
		}

		if errors.Is(err, sqldb.ErrDBIntegrity) || errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return vehiclemodel.VehicleModel{}, vehiclemodel.ErrConflict
		}
		return vehiclemodel.VehicleModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreVehicleModel(dbvm)
}

func (s *Store) Update(ctx context.Context, vm vehiclemodel.VehicleModel) (vehiclemodel.VehicleModel, error) {
	dbvm := toDBVehicleModel(vm)

	const q = `
			UPDATE vehicle_models 
			SET model_series_name 	  = :model_series_name,
				model_commercial_name = :model_commercial_name,
				model_year 			  = :model_year,
				nickname 			  = :nickname,
				engine_displacement   = :engine_displacement,
				drive_type 			  = :drive_type,
				fuel_type 			  = :fuel_type,
				body_style 			  = :body_style,
				transmission_type 	  = :transmission_type
			WHERE vehicle_model_id = :vehicle_model_id
			RETURNING vehicle_model_id
		`
	var pgErr *pgconn.PgError
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbvm, &dbvm); err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "uq_brand_id_series_name_model_year_model_commercial_name" {
				return vehiclemodel.VehicleModel{}, vehiclemodel.ErrDuplicatedModel
			}
		}

		if errors.Is(err, sqldb.ErrDBIntegrity) || errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return vehiclemodel.VehicleModel{}, vehiclemodel.ErrConflict
		}
		return vehiclemodel.VehicleModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreVehicleModel(dbvm)
}

func (s *Store) Delete(ctx context.Context, id int) error {
	data := map[string]any{
		"vehicle_model_id": id,
	}
	const q = `
		DELETE FROM vehicle_models
		WHERE vehicle_model_id = :vehicle_model_id
	`

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return vehiclemodel.ErrNotfound
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
