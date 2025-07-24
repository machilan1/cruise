package seriesvariantdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
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

func (s *Store) NewWithTx(txM tran.TxManager) (seriesvariant.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter seriesvariant.QueryFilter) ([]seriesvariant.SeriesVariant, error) {
	data := map[string]any{}

	const q = `
		SELECT 
			sv.series_variant_id,
			sv.series_variant_name,
			sv.version,
			sv.model_year,
			sv.body_style,
			sv.drive_type,
			sv.fuel_type,
			sv.engine_type,
			sv.engine_displacement,
			sv.valve_count,
			sv.has_turbo,
			sv.transmission_type,
			sv.horse_power,
			json_build_object(
				'brand_series_id',		sv.series_id,
				'brand_series_name',	bs.brand_series_name,
				'brand_id',				b.brand_id,
				'brand_name',			b.brand_name,
				'brand_logo_image',		(SELECT path FROM files WHERE file_id = b.image_id)
			) as series,
			sv.series_id,
			sv.created_at,
			sv.updated_at
		FROM series_variants sv
		LEFT JOIN brand_series bs on sv.series_id = bs.brand_series_id
		LEFT JOIN brands b on bs.brand_id = b.brand_id
	`
	var sb strings.Builder

	sb.WriteString(q)
	applyFilter(filter, data, &sb)

	var dsvs []dbSeriesVariant
	if err := sqldb.NamedQuerySlice(ctx, s.db, sb.String(), data, &dsvs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreSeriesVariants(dsvs)
}

func (s *Store) QueryByID(ctx context.Context, smID int) (seriesvariant.SeriesVariant, error) {
	data := map[string]any{
		"series_variant_id": smID,
	}

	const q = `
		SELECT 
			sv.series_variant_id,
			sv.series_variant_name,
			sv.version,
			sv.model_year,
			sv.body_style,
			sv.drive_type,
			sv.fuel_type,
			sv.engine_type,
			sv.engine_displacement,
			sv.valve_count,
			sv.has_turbo,
			sv.transmission_type,
			sv.horse_power,
			json_build_object(
				'brand_series_id',		sv.series_id,
				'brand_series_name',	bs.brand_series_name,
				'brand_id',				b.brand_id,
				'brand_name',			b.brand_name,
				'brand_logo_image',		(SELECT path FROM files WHERE file_id = b.image_id)
			) as series,
			sv.series_id,
			sv.created_at,
			sv.updated_at
		FROM series_variants sv
		LEFT JOIN brand_series bs on sv.series_id = bs.brand_series_id
		LEFT JOIN brands b on bs.brand_id = b.brand_id
		WHERE sv.series_variant_id  = :series_variant_id
	`

	var dsv dbSeriesVariant
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dsv); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrNotFound
		}
		return seriesvariant.SeriesVariant{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	return toCoreSeriesVariant(dsv)
}

func (s *Store) Create(ctx context.Context, sv seriesvariant.SeriesVariant) (seriesvariant.SeriesVariant, error) {
	const q = `
		INSERT INTO series_variants (
			series_variant_name,
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
			:series_variant_name,
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
			RETURNING series_variant_id
	`

	dsv := toDBSeriesVariant(sv)
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dsv, &dsv); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrConflict
		}
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrConflict
		}
		return seriesvariant.SeriesVariant{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	sv, err := s.QueryByID(ctx, dsv.ID)
	if err != nil {
		return seriesvariant.SeriesVariant{}, fmt.Errorf("querybyid: %w", err)
	}

	return sv, nil
}

func (s *Store) Update(ctx context.Context, sv seriesvariant.SeriesVariant) (seriesvariant.SeriesVariant, error) {
	const q = `
		UPDATE series_variants
		SET series_variant_name		= :series_variant_name,
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
		WHERE series_variant_id 	= :series_variant_id
		RETURNING series_variant_id
	`

	dsv := toDBSeriesVariant(sv)
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dsv, &dsv); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrConflict
		}

		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrNotFound
		}

		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return seriesvariant.SeriesVariant{}, seriesvariant.ErrConflict
		}

		return seriesvariant.SeriesVariant{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	sv, err := s.QueryByID(ctx, dsv.ID)
	if err != nil {
		return seriesvariant.SeriesVariant{}, fmt.Errorf("querybyid: %w", err)
	}

	return sv, nil
}

func (s *Store) Delete(ctx context.Context, svID int) error {
	data := map[string]any{
		"series_variant_id": svID,
	}

	const q = `
		DELETE FROM series_variants WHERE series_variant_id = :series_variant_id
	`

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return seriesvariant.ErrNotFound
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
