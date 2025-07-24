package seriesvariantdb

import (
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
)

func applyFilter(filter seriesvariant.QueryFilter, data map[string]any, sb *strings.Builder) {
	wc := []string{}

	if filter.SeriesID != nil {
		wc = append(wc, "sv.series_id = :series_id")
		data["series_id"] = filter.SeriesID
	}

	if len(wc) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wc, " AND "))
	}
}
