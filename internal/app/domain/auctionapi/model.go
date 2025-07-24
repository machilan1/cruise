package auctionapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/framework/validate"
)

// AppAuction represents an individual auction.
type AppAuction struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	ExecutedCount    int       `json:"executedCount"`
	NotExecutedCount int       `json:"notExecutedCount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	// codegen:{AD}
}

func toAppAuction(auc auction.Auction) AppAuction {
	return AppAuction{
		ID:               auc.ID,
		Name:             auc.Name,
		ExecutedCount:    auc.ExecutedCount,
		NotExecutedCount: auc.NotExecutedCount,
		CreatedAt:        auc.CreatedAt,
		UpdatedAt:        auc.UpdatedAt,
		// codegen:{tAD}
	}
}

func toAppAuctions(aucs []auction.Auction) []AppAuction {
	items := make([]AppAuction, len(aucs))
	for i, auc := range aucs {
		items[i] = toAppAuction(auc)
	}

	return items
}

// =============================================================================

type AppNewAuction struct {
	Name             string `json:"name"`
	ExecutedCount    int    `json:"executedCount"`
	NotExecutedCount int    `json:"notExecutedCount"`
	// codegen:{AN}
}

func (app AppNewAuction) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

func toCoreNewAuction(app AppNewAuction) (auction.NewAuction, error) {
	// codegen:{tManyBN}
	nAuc := auction.NewAuction{
		Name:             app.Name,
		ExecutedCount:    app.ExecutedCount,
		NotExecutedCount: app.NotExecutedCount,
		// codegen:{tBN}
	}

	return nAuc, nil
}

// =============================================================================

type AppUpdateAuction struct {
	Name             *string `json:"name"`
	ExecutedCount    *int    `json:"executedCount"`
	NotExecutedCount *int    `json:"notExecutedCount"`
	// codegen:{AU}
}

func (app AppUpdateAuction) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

func toCoreUpdateAuction(app AppUpdateAuction) (auction.UpdateAuction, error) {
	// codegen:{tManyBU}
	uAuc := auction.UpdateAuction{
		Name:             app.Name,
		ExecutedCount:    app.ExecutedCount,
		NotExecutedCount: app.NotExecutedCount,
		// codegen:{tBU}
	}

	return uAuc, nil
}
