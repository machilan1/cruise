package auction

import "github.com/machilan1/cruise/internal/business/sdk/order"

// DefaultOrderBy is the default order for queries.
var DefaultOrderBy = order.NewBy(OrderByUpdatedAt, order.DESC)

// Set of fields that are allowed to be ordered by.
const (
	OrderByCreatedAt = "created_at"
	OrderByUpdatedAt = "updated_at"
)
