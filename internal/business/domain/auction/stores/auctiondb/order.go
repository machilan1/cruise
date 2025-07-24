package auctiondb

import (
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/sdk/order"
)

var orderByFields = map[string]string{
	auction.OrderByCreatedAt: "created_at",
	auction.OrderByUpdatedAt: "updated_at",
}

func (s *Store) orderByClause(orderBy order.By, sb *strings.Builder) error {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	sb.WriteString(" ORDER BY ")
	sb.WriteString(by)
	sb.WriteString(" ")
	sb.WriteString(orderBy.Direction)
	return nil
}
