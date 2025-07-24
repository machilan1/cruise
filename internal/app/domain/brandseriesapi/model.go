package brandseriesapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/brandseries"
)

type AppBrandSeries struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Brand     AppBrandSeriesBrand `json:"brand"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

func toAppBrandSeries(bs brandseries.BrandSeries) AppBrandSeries {
	return AppBrandSeries{
		ID:   bs.ID,
		Name: bs.Name,
		Brand: AppBrandSeriesBrand{
			ID:        bs.Brand.ID,
			Name:      bs.Brand.Name,
			LogoImage: bs.Brand.LogoImage,
		},
		CreatedAt: bs.CreatedAt,
		UpdatedAt: bs.UpdateAt,
	}
}

func toAppBrandSerieses(bss []brandseries.BrandSeries) []AppBrandSeries {
	abss := make([]AppBrandSeries, len(bss))
	for i, v := range bss {
		abss[i] = toAppBrandSeries(v)
	}
	return abss
}

type AppBrandSeriesBrand struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LogoImage string `json:"logoImage"`
}

type AppNewBrandSeries struct {
	Name    string `json:"name" validate:"require"`
	BrandID int    `json:"brandId" validate:"require"`
}

func toCoreNewBrandSeries(anbs AppNewBrandSeries) (brandseries.NewBrandSeries, error) {
	return brandseries.NewBrandSeries{
		Name:    anbs.Name,
		BrandID: anbs.BrandID,
	}, nil
}

type AppUpdateBrandSeries struct {
	Name    *string `json:"name"`
	BrandID *int    `json:"brandId"`
}

func toCoreUpdateBrandSeries(aubs AppUpdateBrandSeries) (brandseries.UpdateBrandSeries, error) {
	return brandseries.UpdateBrandSeries{
		Name:    aubs.Name,
		BrandID: aubs.BrandID,
	}, nil
}
