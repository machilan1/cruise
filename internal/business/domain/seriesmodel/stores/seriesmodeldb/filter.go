package seriesmodeldb

import (
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
)

func applyFilter(filter seriesmodel.QueryFilter, data map[string]any, sb *strings.Builder) {
	wc := []string{}

	if len(wc) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wc, " AND "))
	}
}
