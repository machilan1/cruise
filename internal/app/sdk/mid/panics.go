package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/machilan1/cruise/internal/app/sdk/metrics"
	"github.com/machilan1/cruise/internal/framework/web"
)

// Panics recovers from panics and converts the panic to an error, so it is
// reported in Metrics and handled in Errors.
func Panics() web.MidFunc {
	m := func(handlerFunc web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return handlerFunc(ctx, w, r)
		}

		return h
	}

	return m
}
