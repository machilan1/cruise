package brandseriesdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/brandseries"
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

func (s *Store) NewWithTx(txM tran.TxManager) (brandseries.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter brandseries.QueryFilter) ([]brandseries.BrandSeries, error) {
	data := map[string]any{}
	const q = `
		SELECT 
			bs.brand_series_id,
			bs.brand_series_name,
			json_build_object(
				'brand_id',		b.brand_id,
				'brand_name',	b.brand_name,
				'logo_image',	b.logo_image
			) as brand,
			bs.created_at,
			bs.updated_at
		FROM brand_series bs
		LEFT JOIN brands b ON bs.brand_id = b.brand_id
	`

	var sb strings.Builder
	sb.WriteString(q)
	applyFilter(filter, data, &sb)

	var dbss []dbBrandSeries
	if err := sqldb.NamedQuerySlice(ctx, s.db, sb.String(), data, &dbss); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreBrandSerieses(dbss)
}

func (s *Store) QueryByID(ctx context.Context, bsID int) (brandseries.BrandSeries, error) {
	data := map[string]any{
		"brand_series_id": bsID,
	}

	const q = `
		SELECT 
			bs.brand_series_id,
			bs.brand_series_name,
			json_build_object(
				'brand_id',		b.brand_id,
				'brand_name',	b.brand_name,
				'logo_image',	b.logo_image
			) as brand,
			bs.created_at,
			bs.updated_at
		FROM brand_series bs
		LEFT JOIN brands b ON bs.brand_id = b.brand_id
		WHERE brand_series_id = :brand_series_id
	`

	var dbs dbBrandSeries
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbs); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return brandseries.BrandSeries{}, brandseries.ErrNotFound
		}
		return brandseries.BrandSeries{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreBrandSeries(dbs)
}

func (s *Store) Create(ctx context.Context, bs brandseries.BrandSeries) (brandseries.BrandSeries, error) {
	dbs := toDBBrandSeries(bs)

	const q = `
		INSERT INTO brand_series(
			brand_series_name,
			brand_id,
			created_at,
			updated_at
		)VALUES(
			:brand_series_name,
			:brand_id,
			:created_at,
			:updated_at
		)
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbs, &dbs); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return brandseries.BrandSeries{}, brandseries.ErrConflict
		}
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return brandseries.BrandSeries{}, brandseries.ErrConflict
		}
		return brandseries.BrandSeries{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	bs, err := s.QueryByID(ctx, bs.ID)
	if err != nil {
		return brandseries.BrandSeries{}, fmt.Errorf("querybyid: %w", err)
	}

	return bs, nil
}

func (s *Store) Update(ctx context.Context, bs brandseries.BrandSeries) (brandseries.BrandSeries, error) {
	dbs := toDBBrandSeries(bs)

	const q = `
		UPDATE brand_series
		SET brand_series_name = :brand_series_name,
			brand_id = :brand_id
		WHERE brand_series_id = :brand_series_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbs, &dbs); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return brandseries.BrandSeries{}, brandseries.ErrConflict
		}
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return brandseries.BrandSeries{}, brandseries.ErrConflict
		}
		return brandseries.BrandSeries{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	bs, err := s.QueryByID(ctx, bs.ID)
	if err != nil {
		return brandseries.BrandSeries{}, fmt.Errorf("querybyid: %w", err)
	}

	return bs, nil
}

func (s *Store) Delete(ctx context.Context, bsID int) error {
	data := map[string]any{
		"brand_series_id": bsID,
	}

	const q = `
		DELETE FROM brand_series WHERE brand_series_id = :brand_series_id
	`
	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return brandseries.ErrNotFound
		}
		return fmt.Errorf("namedquerycontext: %w", err)
	}

	return nil
}
