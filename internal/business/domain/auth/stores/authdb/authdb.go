package authdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

type Store struct {
	db *sqldb.DB
}

func NewStore(db *sqldb.DB) *Store {
	return &Store{
		db: db,
	}
}

// NewWithTx constructs a new Store which replaces the underlying database connection with the provided transaction.
func (s *Store) NewWithTx(txM tran.TxManager) (auth.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) QueryUsers(ctx context.Context) ([]auth.User, error) {
	slog.Warn("authdb.QueryUsers")
	const q = `
		SELECT u.user_id,
			   u.user_type,
			   u.display_name,
		       u.username,
			   u.email,
			   u.created_at,
			   u.updated_at
		FROM users u
		ORDER BY u.created_at, u.user_id
	`

	var dbUsrs []dbUser
	if err := sqldb.NamedQuerySlice(ctx, s.db, q, nil, &dbUsrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	usrs := make([]auth.User, len(dbUsrs))
	for i, dbUsr := range dbUsrs {
		var err error
		usrs[i], err = toCoreUser(dbUsr)
		if err != nil {
			return nil, err
		}
	}

	return usrs, nil
}

// QueryByID retrieves the user with the given ID.
func (s *Store) QueryByID(ctx context.Context, userID int) (auth.User, error) {
	dbUsr := dbUser{
		ID: userID,
	}
	const q = `
		SELECT u.user_id,
			   u.user_type,
			   u.display_name,
		       u.username,
			   u.email,
			   u.created_at,
			   u.updated_at
		FROM users u
		WHERE u.user_id = :user_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbUsr, &dbUsr); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr)
}

func (s *Store) QueryByEmail(ctx context.Context, email string) (auth.User, error) {
	dbUsr := dbUser{
		Email: &email,
	}
	const q = `
		SELECT u.user_id,
			   u.user_type,
			   u.display_name,
		       u.username,
			   u.email,
			   u.created_at,
			   u.updated_at
		FROM users u
		WHERE u.email = :email
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbUsr, &dbUsr); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr)
}

func (s *Store) QueryByUsername(ctx context.Context, username string) (auth.UserWithPassword, error) {
	data := struct {
		Username string `db:"username"`
	}{
		Username: username,
	}

	const q = `
		SELECT u.user_id,
			   u.user_type,
			   u.display_name,
			   u.username,
			   u.email,
			   u.created_at,
			   u.updated_at,
			   up.password
		FROM users u
				 JOIN LATERAL (
			SELECT password
			FROM user_passwords
			WHERE user_id = u.user_id
			ORDER BY created_at DESC
			LIMIT 1
			) up ON TRUE
		WHERE u.username = :username
	`

	var dbUsrWithPW dbUserWithPassword
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbUsrWithPW); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auth.UserWithPassword{}, auth.ErrNotFound
		}
		return auth.UserWithPassword{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUserWithPassword(dbUsrWithPW)
}

func (s *Store) Create(ctx context.Context, usr auth.User) (auth.User, error) {
	dbUsr := toDBUser(usr)

	const q = `
		INSERT INTO users
			(username, display_name, email, created_at, updated_at)
		VALUES (:username, :display_name, :email, :created_at, :updated_at)
		RETURNING user_id
	`
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbUsr, &dbUsr); err != nil {
		// Note: the check here won't work if the underlying database uses a different way to prevent duplicates.
		// For example, if the database uses a unique index(NOT unique constraint), then the error will be different.
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry); err != nil {
			return auth.User{}, auth.ErrUsernameTaken
		}
		return auth.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr)
}

func (s *Store) AddPassword(ctx context.Context, userID int, passwordHash []byte) error {
	slog.Warn("authdb.AddPassword")
	dbUsrPW := toDBPassword(userID, passwordHash)

	const q = `
		INSERT INTO user_passwords
			(user_id, password)
		VALUES (:user_id, :password)
	`
	if err := sqldb.NamedExecContext(ctx, s.db, q, dbUsrPW); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) UpdateUserTypes(ctx context.Context, usrs []auth.User) error {
	// TODO:目前是鎖住不可新增或刪除超級使用者(superAdmin）後續若有需求再打開
	dataU := struct {
		UserIDs   []int    `db:"user_ids"`
		UserTypes []string `db:"user_types"`
	}{
		UserIDs:   make([]int, len(usrs)),
		UserTypes: make([]string, len(usrs)),
	}

	for i, usr := range usrs {
		dataU.UserIDs[i] = usr.ID
		dataU.UserTypes[i] = string(usr.UserType)
	}

	const q = `
		UPDATE users u
		SET user_type = t.user_type
		FROM (
				SELECT * 
				FROM UNNEST(
						CAST(:user_ids AS INT[]),
						CAST(:user_types AS user_types[])
				) 
		) AS t(user_id, user_type)
		WHERE u.user_id = t.user_id
		AND u.user_id NOT IN (SELECT user_id FROM users WHERE user_type = 'superAdmin')
		RETURNING u.user_id
	`
	if err := sqldb.NamedExecContext(ctx, s.db, q, dataU); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) UpdateMe(ctx context.Context, usr auth.User) error {
	dbUsr := toDBUser(usr)

	const q = `
		UPDATE users
		SET display_name = :display_name
		WHERE user_id = :user_id
	`
	if err := sqldb.NamedExecContext(ctx, s.db, q, dbUsr); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) ConsumeAndRefreshEmailQuota(ctx context.Context, email string) error {
	dbUsr := dbUser{
		Email: &email,
	}
	const qReset = `
	WITH uusr AS (
		UPDATE users
		SET email_quota = 3
		WHERE
		CURRENT_DATE> (
		SELECT last_email_refresh_date
		FROM vb_last_email_refresh_date
		)
	)
	UPDATE system_parameters
	SET value = CAST(CURRENT_DATE AS DATE)
	WHERE parameter_key = 'lastEmailRefreshDate'
	AND
	CURRENT_DATE> (
		SELECT last_email_refresh_date
		FROM vb_last_email_refresh_date
		)
	`
	if err := sqldb.NamedExecContext(ctx, s.db, qReset, dbUsr); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	const q = `
		UPDATE users
		SET email_quota = email_quota - 1
		WHERE email = :email
	`
	var pgerr *pgconn.PgError
	if err := sqldb.NamedExecContext(ctx, s.db, q, dbUsr); err != nil {
		if errors.As(err, &pgerr) && pgerr.ConstraintName == "email_quota_check" {
			return auth.ErrorEmailQuota
		}
	}
	return nil
}
