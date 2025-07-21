package brand

import "time"

type Brand struct {
	ID        int
	Name      string
	LogoImage *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewBrand struct {
	Name      string
	LogoImage *string
}

type UpdateBrand struct {
	LogoImage string
}
