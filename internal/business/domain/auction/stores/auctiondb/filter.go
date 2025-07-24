package auctiondb

import (
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/auction"
)

func (s *Store) applyFilter(filter auction.QueryFilter, data map[string]any, sb *strings.Builder) {
	var wc []string

	if filter.Deleted != nil {
		if *filter.Deleted {
			wc = append(wc, "deleted_at IS NOT NULL")
		} else {
			wc = append(wc, "deleted_at IS NULL")
		}
	}

	if len(wc) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wc, " AND "))
	}
}
