package brand

import "time"

type Brand struct {
	ID        int
	Name      string
	Logo      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewBrand struct {
	Name string
	Logo *string
}

type UpdateBrand struct {
	Logo string
}
