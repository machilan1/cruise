package brandseriesdb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/brandseries"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbBrandSeries struct {
	ID        int                                   `db:"brand_series_id"`
	Name      string                                `db:"brand_series_name"`
	Brand     dbjson.JSONColumn[dbBrandSeriesBrand] `db:"brand"`
	BrandID   int                                   `db:"brand_id"`
	CreatedAt time.Time                             `db:"created_at"`
	UpdatedAt time.Time                             `db:"updated_at"`
}

type dbBrandSeriesBrand struct {
	ID        int    `json:"brand_id"`
	Name      string `json:"brand_name"`
	LogoImage string `json:"logo_image"`
}

func toDBBrandSeries(bs brandseries.BrandSeries) dbBrandSeries {
	return dbBrandSeries{
		ID:        bs.ID,
		Name:      bs.Name,
		BrandID:   bs.Brand.ID,
		CreatedAt: bs.CreatedAt,
		UpdatedAt: bs.UpdateAt,
	}
}

func toCoreBrandSeries(dbs dbBrandSeries) (brandseries.BrandSeries, error) {
	return brandseries.BrandSeries{
		ID:   dbs.ID,
		Name: dbs.Name,
		Brand: brandseries.BrandSeriesBrand{
			ID:        dbs.Brand.Get().ID,
			Name:      dbs.Brand.Get().Name,
			LogoImage: dbs.Brand.Get().LogoImage,
		},
		CreatedAt: dbs.CreatedAt,
		UpdateAt:  dbs.UpdatedAt,
	}, nil
}

func toCoreBrandSerieses(dbss []dbBrandSeries) ([]brandseries.BrandSeries, error) {
	bss := make([]brandseries.BrandSeries, len(dbss))
	for i, v := range dbss {
		bs, err := toCoreBrandSeries(v)
		if err != nil {
			return nil, err
		}
		bss[i] = bs
	}
	return bss, nil
}
