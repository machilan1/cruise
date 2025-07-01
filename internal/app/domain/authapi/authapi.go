package authapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/notification"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

var (
	ErrInvalidToken = errs.NewTrustedError(errors.New("invalid token"), http.StatusBadRequest)
	ErrEmailQuota   = errs.NewTrustedError(errors.New("email quota exceeded"), http.StatusConflict)
)

type handlers struct {
	log   *logger.Logger
	txM   tran.TxManager
	sess  *sess.Manager
	auth  *auth.Core
	notif *notification.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, sess *sess.Manager, auth *auth.Core, notif *notification.Core) *handlers {
	return &handlers{
		log:   log,
		txM:   txM,
		sess:  sess,
		auth:  auth,
		notif: notif,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	authCore, err := h.auth.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:   h.log,
		txM:   txM,
		sess:  h.sess,
		auth:  authCore,
		notif: h.notif,
	}, nil
}

func (h *handlers) register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppRegisterInput
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}
	nUsr, err := toCoreNewUser(app)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	var usr auth.User
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		usr, err = h.auth.Register(ctx, nUsr)
		if err != nil {
			if errors.Is(err, auth.ErrUsernameTaken) {
				return errs.NewTrustedError(errors.New("username taken"), http.StatusConflict)
			}
			return fmt.Errorf("register: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	// TODO: send welcome email

	return web.Respond(ctx, w, toAppMe(usr), http.StatusCreated)
}

func (h *handlers) login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppLoginInput
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	usr, err := h.auth.Authenticate(ctx, app.Username, app.Password)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			h.log.Info(ctx, "login: no user found in the system", "username", app.Username)
			return errs.NewTrustedError(errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		}
		if errors.Is(err, auth.ErrAuthFailed) {
			h.log.Info(ctx, "login: failed authentication attempt", "username", app.Username)
			return errs.NewTrustedError(errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		}
		return fmt.Errorf("authenticate: %w", err)
	}

	session, err := h.sess.NewContext(ctx, r, mid.AuthSessionID)
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	if !session.IsNew {
		if err := h.sess.Delete(ctx, w, r, session); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
		session, err = h.sess.NewContext(ctx, r, mid.AuthSessionID)
		if err != nil {
			return fmt.Errorf("new session: %w", err)
		}
	}

	session.Values[sess.SessionUserID] = usr.ID
	if err := h.sess.SaveContext(ctx, w, r, session); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	return web.Respond(ctx, w, toAppMe(usr), http.StatusOK)
}

func (h *handlers) logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	session, err := h.sess.Get(r, mid.AuthSessionID)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	if err := h.sess.Delete(ctx, w, r, session); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) getMe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	return web.Respond(ctx, w, toAppMe(usr), http.StatusOK)
}

func (h *handlers) resetPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppResetPasswordInput
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	nPwd, err := auth.ParsePassword(app.NewPassword)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	uID := web.Param(r, "userID")
	userID, err := strconv.Atoi(uID)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err = h.auth.ResetPassword(ctx, nPwd, userID); err != nil {
			return fmt.Errorf("update authz: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) forgotPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppForgotPasswordInput
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	email, err := mail.ParseAddress(app.Email)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	token, err := h.auth.RequestPasswordResetToken(ctx, *email)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			h.log.Info(ctx, "forgot password: no user found in the system", "email", email)

			// Do not leak information about the user's existence.
			return web.Respond(ctx, w, nil, http.StatusNoContent)
		}
		return fmt.Errorf("request password reset token: %w", err)
	}

	// check and consume quota before sending email
	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.auth.ConsumeAndRefreshEmailQuota(ctx, *email); err != nil {
			if errors.Is(err, auth.ErrorEmailQuota) {
				return ErrEmailQuota
			}
			return fmt.Errorf("consume and refresh email quota: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	if err := h.notif.SendPasswordResetEmail(ctx, *email, token); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) resetPasswordWithToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppResetPasswordWithTokenInput
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	pw, err := auth.ParsePassword(app.NewPassword)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.auth.ResetPasswordWithToken(ctx, app.Token, pw); err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				return ErrInvalidToken
			}
			return fmt.Errorf("reset password with token: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) updateUserTypes(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateUserTypeInputs
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	uUsers, err := toCoreUpdateUserTypeInputs(app)
	if err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}
	slog.Warn("msg", "ID", uUsers[0].UserID)

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if err := h.auth.UpdateUserTypes(ctx, uUsers); err != nil {
			return fmt.Errorf("update user type fail: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *handlers) queryUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	usrs, err := h.auth.QueryUsers(ctx)
	if err != nil {
		return fmt.Errorf("query users: %w", err)
	}

	return web.Respond(ctx, w, toAppUsers(usrs), http.StatusOK)
}

func (h *handlers) updateMe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	var app AppUpdateMe
	if err := web.Decode(r, &app); err != nil {
		return errs.NewTrustedError(err, http.StatusBadRequest)
	}

	uMe := toCoreUpdateMe(app)

	var result auth.User

	if err := h.txM.RunTx(ctx, func(txM tran.TxManager) error {
		h, err := h.newWithTx(txM)
		if err != nil {
			return err
		}

		if result, err = h.auth.UpdateMe(ctx, usr.ID, uMe); err != nil {
			return fmt.Errorf("update user type fail: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppMe(result), http.StatusOK)
}
