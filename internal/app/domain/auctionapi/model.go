package auctionapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/framework/validate"
)

// AppAuction represents an individual auction.
type AppAuction struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// codegen:{AD}
}

func toAppAuction(auc auction.Auction) AppAuction {
	return AppAuction{
		ID:        auc.ID,
		Date:      auc.Date,
		Note:      auc.Note,
		CreatedAt: auc.CreatedAt,
		UpdatedAt: auc.UpdatedAt,
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
	Date time.Time `json:"date" validate:"required"`
	Note *string   `json:"note"`
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
	var note string
	if app.Note != nil {
		note = *app.Note
	}

	nAuc := auction.NewAuction{
		Date: app.Date,
		Note: note,
		// codegen:{tBN}
	}

	return nAuc, nil
}

// =============================================================================

type AppUpdateAuction struct {
	Date *time.Time `json:"date"`
	Note *string    `json:"note"`
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
		Date: app.Date,
		Note: app.Note,
		// codegen:{tBU}
	}

	return uAuc, nil
}
