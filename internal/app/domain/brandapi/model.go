package brandapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/brand"
)

type AppBrand struct {
	ID        int       `json:"brandId"`
	Name      string    `json:"name"`
	LogoImage *string   `json:"logoImage"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toAppBrand(brd brand.Brand) AppBrand {
	return AppBrand{
		ID:        brd.ID,
		Name:      brd.Name,
		LogoImage: brd.LogoImage,
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
	Name      string  `json:"name"`
	LogoImage *string `json:"logoImage"`
}

func toCoreNewBrand(abrd AppNewBrand) brand.NewBrand {
	return brand.NewBrand{
		Name:      abrd.Name,
		LogoImage: abrd.LogoImage,
	}
}

type AppUpdateBrand struct {
	LogoImage string `json:"logoImage"`
}

func toCoreUpdateBrand(aubrd AppUpdateBrand) brand.UpdateBrand {
	return brand.UpdateBrand{
		LogoImage: aubrd.LogoImage,
	}
}
