package auction

import "time"

type Auction struct {
	ID        int
	Date      time.Time
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// codegen:{BD}
}

type NewAuction struct {
	Date time.Time
	Note string
	// codegen:{BN}
}

type UpdateAuction struct {
	Date *time.Time
	Note *string
	// codegen:{BU}
}
