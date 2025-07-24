package auctiondb

import (
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auction"
)

// dbAuction represents an individual auction.
type dbAuction struct {
	ID               int        `db:"admin_transaction_id"`
	Name             string     `db:"name"`
	ExecutedCount    int        `db:"executed_count"`
	NotExecutedCount int        `db:"not_executed_count"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	// codegen:{SD}
}

func toDBAuction(auc auction.Auction) dbAuction {
	dbAuc := dbAuction{
		ID:               auc.ID,
		Name:             auc.Name,
		ExecutedCount:    auc.ExecutedCount,
		NotExecutedCount: auc.NotExecutedCount,
		CreatedAt:        auc.CreatedAt,
		UpdatedAt:        auc.UpdatedAt,
		DeletedAt:        auc.DeletedAt,
		// codegen:{tSD}
	}

	return dbAuc
}

func toCoreAuction(dbAuc dbAuction) (auction.Auction, error) {
	auc := auction.Auction{
		ID:               dbAuc.ID,
		Name:             dbAuc.Name,
		ExecutedCount:    dbAuc.ExecutedCount,
		NotExecutedCount: dbAuc.NotExecutedCount,
		CreatedAt:        dbAuc.CreatedAt,
		UpdatedAt:        dbAuc.UpdatedAt,
		DeletedAt:        dbAuc.DeletedAt,
		// codegen:{tBD}
	}

	return auc, nil
}

func toCoreAuctions(dbAuctions []dbAuction) ([]auction.Auction, error) {
	aucs := make([]auction.Auction, len(dbAuctions))
	for i, dbAuc := range dbAuctions {
		auc, err := toCoreAuction(dbAuc)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
		aucs[i] = auc
	}

	return aucs, nil
}
