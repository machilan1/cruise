package authapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/business/domain/notification"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type Config struct {
	Log          *logger.Logger
	TxM          tran.TxManager
	Sess         *sess.Manager
	Auth         *auth.Core
	Notification *notification.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authn := mid.AuthN(cfg.Sess, cfg.Auth)
	authSuper := mid.AuthNAndIsSuperAdmin(cfg.Sess, cfg.Auth)
	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.Sess, cfg.Auth, cfg.Notification)

	authApp := *app
	authApp.PostMid = []web.MidFunc{
		authn,
	}
	authSuperAdmin := *app
	authSuperAdmin.PostMid = []web.MidFunc{
		authSuper,
	}
	// --- Public ---
	app.HandleFunc(http.MethodPost, version, "/auth/register", hdl.register)
	app.HandleFunc(http.MethodPost, version, "/auth/login", hdl.login)
	app.HandleFunc(http.MethodPost, version, "/auth/forgot-password", hdl.forgotPassword)
	app.HandleFunc(http.MethodPost, version, "/auth/reset-password", hdl.resetPasswordWithToken)
	// --- Requires authentication ---
	authApp.HandleFunc(http.MethodGet, version, "/auth/me", hdl.getMe)
	authApp.HandleFunc(http.MethodPut, version, "/auth/me", hdl.updateMe)
	authApp.HandleFunc(http.MethodPost, version, "/auth/logout", hdl.logout)

	authSuperAdmin.HandleFunc(http.MethodGet, version, "/auth/users", hdl.queryUsers)
	// TODO:目前是鎖住不可新增或刪除超級使用者(superAdmin）後續若有需求再打開
	authSuperAdmin.HandleFunc(http.MethodPut, version, "/auth/reset-user-types", hdl.updateUserTypes)
	authSuperAdmin.HandleFunc(http.MethodPut, version, "/auth/users/{userID}/reset-password", hdl.resetPassword)
}
