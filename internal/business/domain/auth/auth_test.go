package auth_test

import (
	"context"
	"errors"
	"net/mail"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/auth/stores/authdb"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/testhelper"
)

var testDatabaseInstance *sqldb.TestInstance

func TestMain(m *testing.M) {
	testDatabaseInstance = sqldb.MustTestInstance()
	defer testDatabaseInstance.MustClose()
	m.Run()
}

type testSuite struct {
	auth *auth.Core
}

func newTestSuite(t *testing.T) *testSuite {
	t.Helper()

	log := testhelper.TestLogger(t)
	testDB, _ := testDatabaseInstance.NewDatabase(t, log)

	return &testSuite{
		auth: auth.NewCore(authdb.NewStore(testDB), []byte("secret")),
	}
}

func TestAuth_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	want := auth.User{
		ID:        1,
		Username:  auth.MustParseUsername("johndoe"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pwPlain := "password123"
	pw := auth.MustParsePassword(pwPlain)

	// Register a new user.
	{
		nUsr := toNewUser(want, pw)
		got, err := ts.auth.Register(ctx, nUsr)
		if err != nil {
			t.Fatal(err)
		}
		checkUser(t, got, want)
	}

	// Query the user by ID.
	{
		got, err := ts.auth.QueryByID(ctx, want.ID)
		if err != nil {
			t.Fatal(err)
		}
		checkUser(t, got, want)
	}

	// Query the user by username.
	{
		got, err := ts.auth.QueryByUsername(ctx, want.Username.String())
		if err != nil {
			t.Fatal(err)
		}
		checkUser(t, got, want)
	}

	// Login with the user.
	{
		got, err := ts.auth.Authenticate(ctx, want.Username.String(), pwPlain)
		if err != nil {
			t.Fatal(err)
		}
		checkUser(t, got, want)
	}
}

func TestAuth_QueryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	if _, err := ts.auth.QueryByID(ctx, 999); !errors.Is(err, auth.ErrNotFound) {
		t.Errorf("got err: %v, want: %v", err, auth.ErrNotFound)
	}
	if _, err := ts.auth.QueryByUsername(ctx, "unknown"); !errors.Is(err, auth.ErrNotFound) {
		t.Errorf("got err: %v, want: %v", err, auth.ErrNotFound)
	}
}

func TestAuth_RegisterError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	usr, _ := mustRegisterUser(ctx, t, ts.auth)

	validNewUser := auth.NewUser{
		Username: auth.MustParseUsername("validusername"),
		Password: auth.MustParsePassword("validpassword123"),
	}

	tests := []struct {
		name     string
		username *auth.Username
		email    *mail.Address
		err      error
	}{
		{"taken username", &usr.Username, nil, auth.ErrUsernameTaken},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nu := validNewUser
			if tt.username != nil {
				nu.Username = *tt.username
			}

			if _, err := ts.auth.Register(ctx, nu); !errors.Is(err, tt.err) {
				t.Errorf("got err: %v, want: %v", err, tt.err)
			}
		})
	}
}

func TestAuth_AuthenticateError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	usr, pw := mustRegisterUser(ctx, t, ts.auth)

	tests := []struct {
		name     string
		username string
		pw       string
	}{
		{"empty username and password", "", ""},
		{"empty username", "", pw.String()},
		{"empty password", usr.Username.String(), ""},
		{"wrong password", usr.Username.String(), "wrongpassword"},
		{"non-existent username", "nonexistent", pw.String()},
		{"invalid username format", "不正なユーザー名", pw.String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := ts.auth.Authenticate(ctx, tt.username, tt.pw); !errors.Is(err, auth.ErrAuthFailed) {
				t.Errorf("got err: %v, want: %v", err, auth.ErrAuthFailed)
			}
		})
	}
}

func checkUser(t *testing.T, got, want auth.User) {
	t.Helper()

	if diff := cmp.Diff(got, want, sqldb.ApproxTime); diff != "" {
		t.Fatalf("mismatch (-got +want):\n%s", diff)
	}
}

func toNewUser(u auth.User, pw auth.Password) auth.NewUser {
	return auth.NewUser{
		Username: u.Username,
		Password: pw,
	}
}

func mustRegisterUser(ctx context.Context, t *testing.T, c *auth.Core) (auth.User, auth.Password) {
	t.Helper()

	pw := auth.MustParsePassword("password123")

	usr, err := c.Register(ctx, auth.NewUser{
		Username: auth.MustParseUsername("johndoe"),
		Password: pw,
	})
	if err != nil {
		t.Fatal(err)
	}
	return usr, pw
}
