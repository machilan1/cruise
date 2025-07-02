package vehiclemodeldb

import (
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
)

func applyFilter(filter vehiclemodel.QueryFilter, data map[string]any, sb *strings.Builder) {
	var wc []string

	if filter.BrandID != nil {
		wc = append(wc, "brand_id = :brand_id")
	}

	if len(wc) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wc, " AND "))
	}
}
