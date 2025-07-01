package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/framework/web"
)

// DevMode is a flag that can be set to true to bypass the session validation
// and use a hardcoded user ID. This is useful for development and testing.
//
// For instance, when DevMode is true, one can send a request with a cookie like:
//
//	Cookie: _session=1
//
// and the request will be authenticated as user with ID 1.
var DevMode = false

const (
	// AuthSessionID is the session RoleID used for the auth session.
	// This is used as the key for saving and retrieving the session in cookie.
	AuthSessionID = "_session"
)

var errUnauthorized = errs.NewTrustedError(errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)

func setUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int, error) {
	v, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	return v, nil
}

func setUser(ctx context.Context, usr auth.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

func GetUser(ctx context.Context) (auth.User, error) {
	usr, ok := ctx.Value(userKey).(auth.User)
	if !ok {
		return auth.User{}, errors.New("user not found in context")
	}
	return usr, nil
}

func AuthN(sm *sess.Manager, authCore *auth.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID int
			if DevMode {
				c, err := r.Cookie(AuthSessionID)
				if err != nil {
					return errUnauthorized
				}
				userID, err = strconv.Atoi(c.Value)
				if err != nil {
					return errUnauthorized
				}
			} else {
				session, err := sm.GetContext(ctx, r, AuthSessionID)
				if err != nil {
					return fmt.Errorf("session get: id[%s]: %w", AuthSessionID, err)
				}

				var ok bool
				userID, ok = session.Values[sess.SessionUserID].(int)
				if !ok {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					// TODO: log this error to a monitoring system.
					_ = sm.Delete(ctx, w, r, session)
					return errUnauthorized
				}
			}
			ctx = setUserID(ctx, userID)

			usr, err := authCore.QueryByID(ctx, userID)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					return errUnauthorized
				}
				return fmt.Errorf("query: id[%d]: %w", userID, err)
			}
			ctx = setUser(ctx, usr)

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func AuthNAndIsSuperAdmin(sm *sess.Manager, authCore *auth.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID int
			if DevMode {
				c, err := r.Cookie(AuthSessionID)
				if err != nil {
					return errUnauthorized
				}
				userID, err = strconv.Atoi(c.Value)
				if err != nil {
					return errUnauthorized
				}
			} else {
				session, err := sm.GetContext(ctx, r, AuthSessionID)
				if err != nil {
					return fmt.Errorf("session get: id[%s]: %w", AuthSessionID, err)
				}

				var ok bool
				userID, ok = session.Values[sess.SessionUserID].(int)
				if !ok {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					// TODO: log this error to a monitoring system.
					_ = sm.Delete(ctx, w, r, session)
					return errUnauthorized
				}
			}
			ctx = setUserID(ctx, userID)

			usr, err := authCore.QueryByID(ctx, userID)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					return errUnauthorized
				}
				return fmt.Errorf("query: id[%d]: %w", userID, err)
			}
			if !(usr.UserType == auth.UserTypeSuperAdmin) {
				return errUnauthorized
			}
			ctx = setUser(ctx, usr)

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func AuthNAndIsAdmin(sm *sess.Manager, authCore *auth.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID int
			if DevMode {
				c, err := r.Cookie(AuthSessionID)
				if err != nil {
					return errUnauthorized
				}
				userID, err = strconv.Atoi(c.Value)
				if err != nil {
					return errUnauthorized
				}
			} else {
				session, err := sm.GetContext(ctx, r, AuthSessionID)
				if err != nil {
					return fmt.Errorf("session get: id[%s]: %w", AuthSessionID, err)
				}

				var ok bool
				userID, ok = session.Values[sess.SessionUserID].(int)
				if !ok {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					// TODO: log this error to a monitoring system.
					_ = sm.Delete(ctx, w, r, session)
					return errUnauthorized
				}
			}
			ctx = setUserID(ctx, userID)

			usr, err := authCore.QueryByID(ctx, userID)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					return errUnauthorized
				}
				return fmt.Errorf("query: id[%d]: %w", userID, err)
			}
			if !(usr.UserType == auth.UserTypeAdmin) {
				return errUnauthorized
			}
			ctx = setUser(ctx, usr)

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func AuthNAndIsOneOfAdmin(sm *sess.Manager, authCore *auth.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID int
			if DevMode {
				c, err := r.Cookie(AuthSessionID)
				if err != nil {
					return errUnauthorized
				}
				userID, err = strconv.Atoi(c.Value)
				if err != nil {
					return errUnauthorized
				}
			} else {
				session, err := sm.GetContext(ctx, r, AuthSessionID)
				if err != nil {
					return fmt.Errorf("session get: id[%s]: %w", AuthSessionID, err)
				}

				var ok bool
				userID, ok = session.Values[sess.SessionUserID].(int)
				if !ok {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					// TODO: log this error to a monitoring system.
					_ = sm.Delete(ctx, w, r, session)
					return errUnauthorized
				}
			}
			ctx = setUserID(ctx, userID)

			usr, err := authCore.QueryByID(ctx, userID)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					return errUnauthorized
				}
				return fmt.Errorf("query: id[%d]: %w", userID, err)
			}
			if !(usr.UserType == auth.UserTypeSuperAdmin || usr.UserType == auth.UserTypeAdmin) {
				return errUnauthorized
			}
			ctx = setUser(ctx, usr)

			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func AuthNAndStaffOrOneOfAdmin(sm *sess.Manager, authCore *auth.Core) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID int
			if DevMode {
				c, err := r.Cookie(AuthSessionID)
				if err != nil {
					return errUnauthorized
				}
				userID, err = strconv.Atoi(c.Value)
				if err != nil {
					return errUnauthorized
				}
			} else {
				session, err := sm.GetContext(ctx, r, AuthSessionID)
				if err != nil {
					return fmt.Errorf("session get: id[%s]: %w", AuthSessionID, err)
				}

				var ok bool
				userID, ok = session.Values[sess.SessionUserID].(int)
				if !ok {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					// TODO: log this error to a monitoring system.
					_ = sm.Delete(ctx, w, r, session)
					return errUnauthorized
				}
			}
			ctx = setUserID(ctx, userID)

			usr, err := authCore.QueryByID(ctx, userID)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					// We are ignoring the error here since the session is invalid,
					// and we can safely return an unauthorized error instead.
					return errUnauthorized
				}
				return fmt.Errorf("query: id[%d]: %w", userID, err)
			}
			if !(usr.UserType == auth.UserTypeSuperAdmin || usr.UserType == auth.UserTypeAdmin || usr.UserType == auth.UserTypeStaff) {
				return errUnauthorized
			}
			ctx = setUser(ctx, usr)

			return next(ctx, w, r)
		}

		return h
	}

	return m
}
