package brandapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/brand"
)

type AppBrand struct {
	ID        int       `json:"brandId"`
	Name      string    `json:"name"`
	Logo      *string   `json:"logo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toAppBrand(brd brand.Brand) AppBrand {
	return AppBrand{
		ID:        brd.ID,
		Name:      brd.Name,
		Logo:      brd.Logo,
		CreatedAt: brd.CreatedAt,
		UpdatedAt: brd.UpdatedAt,
	}
}

func toAppBrands(brds []brand.Brand) []AppBrand {
	abrds := make([]AppBrand, len(brds))
	for i, v := range brds {
		abrds[i] = toAppBrand(v)
	}
	return abrds
}

type AppNewBrand struct {
	Name string  `json:"name"`
	Logo *string `json:"logo"`
}

func toCoreNewBrand(abrd AppNewBrand) brand.NewBrand {
	return brand.NewBrand{
		Name: abrd.Name,
		Logo: abrd.Logo,
	}
}

type AppUpdateBrand struct {
	Logo string `json:"logo"`
}

func toCoreUpdateBrand(aubrd AppUpdateBrand) brand.UpdateBrand {
	return brand.UpdateBrand{
		Logo: aubrd.Logo,
	}
}
