package brandseries

import "time"

type BrandSeries struct {
	ID        int
	Name      string
	Brand     BrandSeriesBrand
	CreatedAt time.Time
	UpdateAt  time.Time
}

type BrandSeriesBrand struct {
	ID        int
	Name      string
	LogoImage string
}

type NewBrandSeries struct {
	Name    string
	BrandID int
}

type UpdateBrandSeries struct {
	Name    *string
	BrandID *int
}
