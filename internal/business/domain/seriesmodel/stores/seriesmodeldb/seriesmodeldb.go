package seriesmodeldb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
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

func (s *Store) NewWithTx(txM tran.TxManager) (seriesmodel.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter seriesmodel.QueryFilter) ([]seriesmodel.SeriesModel, error) {
	data := map[string]any{}

	const q = `
		SELECT 
			sm.series_model_id,
			sm.series_model_name,
			sm.version,
			sm.model_year,
			sm.body_style,
			sm.drive_type,
			sm.fuel_type,
			sm.engine_type,
			sm.engine_displacement,
			sm.valve_count,
			sm.has_turbo,
			sm.transmission_type,
			sm.horse_power,
			json_build_object(
				'brand_series_id',		sm.series_id,
				'brand_series_name',	bs.brand_series_name,
				'brand_id',				b.brand_id,
				'brand_name',			b.brand_name,
				'brand_logo_image',		(SELECT path FROM files WHERE file_id = b.image_id)
			) as series,
			sm.series_id,
			sm.created_at,
			sm.updated_at
		FROM series_models sm
		LEFT JOIN brand_series bs on sm.series_id = bs.brand_series_id
		LEFT JOIN brands b on bs.brand_id = b.brand_id
	`
	var sb strings.Builder

	sb.WriteString(q)
	applyFilter(filter, data, &sb)

	var dsms []dbSeriesModel
	if err := sqldb.NamedQuerySlice(ctx, s.db, sb.String(), data, &dsms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreSeriesModels(dsms)
}

func (s *Store) QueryByID(ctx context.Context, smID int) (seriesmodel.SeriesModel, error) {
	data := map[string]any{
		"series_model_id": smID,
	}

	const q = `
		SELECT 
			sm.series_model_id,
			sm.series_model_name,
			sm.version,
			sm.model_year,
			sm.body_style,
			sm.drive_type,
			sm.fuel_type,
			sm.engine_type,
			sm.engine_displacement,
			sm.valve_count,
			sm.has_turbo,
			sm.transmission_type,
			sm.horse_power,
			json_build_object(
				'brand_series_id',		sm.series_id,
				'brand_series_name',	bs.brand_series_name,
				'brand_id',				b.brand_id,
				'brand_name',			b.brand_name,
				'brand_logo_image',		(SELECT path FROM files WHERE file_id = b.image_id)
			) as series,
			sm.series_id,
			sm.created_at,
			sm.updated_at
		FROM series_models sm
		LEFT JOIN brand_series bs on sm.series_id = bs.brand_series_id
		LEFT JOIN brands b on bs.brand_id = b.brand_id
		WHERE sm.series_model_id  = :series_model_id
	`

	var dsm dbSeriesModel
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dsm); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrNotFound
		}
		return seriesmodel.SeriesModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	return toCoreSeriesModel(dsm)
}

func (s *Store) Create(ctx context.Context, sm seriesmodel.SeriesModel) (seriesmodel.SeriesModel, error) {
	const q = `
		INSERT INTO series_models (
			series_model_name,
			version,
			model_year,
			body_style,
			drive_type,
			fuel_type,
			engine_type,
			engine_displacement,
			valve_count,
			has_turbo,
			transmission_type,
			horse_power,
			series_id,
			created_at,
			updated_at
		)VALUES(
			:series_model_name,
			:version,
			:model_year,
			:body_style,
			:drive_type,
			:fuel_type,
			:engine_type,
			:engine_displacement,
			:valve_count,
			:has_turbo,
			:transmission_type,
			:horse_power,
			:series_id,
			:created_at,
			:updated_at
		)
			RETURNING series_model_id
	`

	dsm := toDBSeriesModels(sm)
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dsm, &dsm); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrConflict
		}
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrConflict
		}
		return seriesmodel.SeriesModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	sm, err := s.QueryByID(ctx, dsm.ID)
	if err != nil {
		return seriesmodel.SeriesModel{}, fmt.Errorf("querybyid: %w", err)
	}

	return sm, nil
}

func (s *Store) Update(ctx context.Context, sm seriesmodel.SeriesModel) (seriesmodel.SeriesModel, error) {
	const q = `
		UPDATE series_models
		SET series_model_name		= :series_model_name,
			version					= :version,
			model_year				= :model_year,
			body_style				= :body_style,
			drive_type				= :drive_type,
			fuel_type				= :fuel_type,
			engine_type				= :engine_type,
			engine_displacement		= :engine_displacement,
			valve_count				= :valve_count,
			has_turbo				= :has_turbo,
			transmission_type		= :transmission_type,
			horse_power				= :horse_power,
			series_id				= :series_id,
			updated_at				= :updated_at
		WHERE series_model_id 	= :series_model_id
		RETURNING series_model_id
	`

	dsm := toDBSeriesModels(sm)
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dsm, &dsm); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrConflict
		}

		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrNotFound
		}

		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return seriesmodel.SeriesModel{}, seriesmodel.ErrConflict
		}

		return seriesmodel.SeriesModel{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	sm, err := s.QueryByID(ctx, dsm.ID)
	if err != nil {
		return seriesmodel.SeriesModel{}, fmt.Errorf("querybyid: %w", err)
	}

	return sm, nil
}

func (s *Store) Delete(ctx context.Context, smID int) error {
	data := map[string]any{
		"series_model_id": smID,
	}

	const q = `
		DELETE FROM series_models WHERE series_model_id = :series_model_id
	`

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesmodel.ErrNotFound
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
