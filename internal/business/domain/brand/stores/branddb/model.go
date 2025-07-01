package branddb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/brand"
)

type dbBrand struct {
	ID        int       `db:"brand_id"`
	Name      string    `db:"brand_name"`
	Logo      *string   `db:"logo"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func toCoreBrand(dbb dbBrand) brand.Brand {
	return brand.Brand{
		ID:   dbb.ID,
		Name: dbb.Name,
		Logo: dbb.Logo,
	}
}

func toCoreBrands(dbbs []dbBrand) []brand.Brand {
	bs := make([]brand.Brand, len(dbbs))
	for i, v := range dbbs {
		bs[i] = toCoreBrand(v)
	}
	return bs
}

func toDBBrand(brd brand.Brand) dbBrand {
	return dbBrand{
		ID:   brd.ID,
		Name: brd.Name,
		Logo: brd.Logo,
	}
}
