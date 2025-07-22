package auctionhousedb

import (
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
)

func applyFilter(filter auctionhouse.QueryFilter, data map[string]any, sb *strings.Builder) {
	var wc []string

	if len(wc) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wc, " AND "))
	}
}
