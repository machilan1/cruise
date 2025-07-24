package auction

import "time"

type Auction struct {
	ID               int
	Name             string
	ExecutedCount    int
	NotExecutedCount int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time

	// codegen:{BD}
}

type NewAuction struct {
	Name             string
	ExecutedCount    int
	NotExecutedCount int
	// codegen:{BN}
}

type UpdateAuction struct {
	Name             *string
	ExecutedCount    *int
	NotExecutedCount *int
	// codegen:{BU}
}
