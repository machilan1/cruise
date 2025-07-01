package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strconv"
	"time"

	// Load the timezone data
	// Since when deploying to the container, the timezone data is not included in the image.
	// We need to import the tzdata package to include the timezone data.
	_ "time/tzdata"

	"github.com/golang-jwt/jwt/v5"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound        = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username is taken")
	ErrAuthFailed      = errors.New("authentication failed")
	ErrDuplicateUserID = errors.New("duplicate userID")
	ErrInvalidToken    = errors.New("invalid token")
	ErrorEmailQuota    = errors.New("email quota exceeded")
)

var pwResetExpiresIn = 1 * time.Hour

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	QueryUsers(ctx context.Context) ([]User, error)
	QueryByID(ctx context.Context, userID int) (User, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
	QueryByUsername(ctx context.Context, username string) (UserWithPassword, error)
	Create(ctx context.Context, usr User) (User, error)
	AddPassword(ctx context.Context, userID int, passwordHash []byte) error
	UpdateUserTypes(ctx context.Context, usrs []User) error
	UpdateMe(ctx context.Context, usr User) error

	ConsumeAndRefreshEmailQuota(ctx context.Context, email string) error
}

// ====================================================================================

type Core struct {
	storer Storer
	// secret is used to sign JWT tokens.
	secret []byte
}

func NewCore(storer Storer, secret []byte) *Core {
	return &Core{
		storer: storer,
		secret: secret,
	}
}

func (c *Core) NewWithTx(txM tran.TxManager) (*Core, error) {
	storer, err := c.storer.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &Core{
		storer: storer,
		secret: c.secret,
	}, nil
}

// QueryByID retrieves the user with the given RoleID.
// It returns ErrNotFound if the user is not found.
func (c *Core) QueryByID(ctx context.Context, userID int) (User, error) {
	usr, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%d]: %w", userID, err)
	}

	return usr, nil
}

// QueryByUsername retrieves the user with the given username.
// It returns ErrNotFound if the user is not found.
func (c *Core) QueryByUsername(ctx context.Context, username string) (User, error) {
	usr, err := c.storer.QueryByUsername(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("query: username[%s]: %w", username, err)
	}
	return usr.User, nil
}

// Register creates a new user and stores it in the database.
// The password is hashed before being stored.
// It returns ErrUsernameTaken if the username is already taken.
func (c *Core) Register(ctx context.Context, nUsr NewUser) (User, error) {
	now := time.Now()
	usr := User{
		Username:    nUsr.Username,
		DisplayName: nUsr.DisplayName,
		UserType:    UserTypeUser,
		Email:       nUsr.Email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	usr, err := c.storer.Create(ctx, usr)
	if err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	pwHash, err := nUsr.Password.Hash()
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}
	if err := c.storer.AddPassword(ctx, usr.ID, pwHash); err != nil {
		return User{}, fmt.Errorf("add password: %w", err)
	}

	result, err := c.storer.QueryByID(ctx, usr.ID)
	if err != nil {
		return User{}, fmt.Errorf("query user after create: %w", err)
	}

	return result, nil
}

// Authenticate checks if the provided username or email and password are valid.
// We do NOT ask for a concrete `Password` type here, since there might be some cases where the password is already
// saved, and yet we updated the parsing rules for the `Password` type. In this case, the old password would not be
// parsable anymore, and we would not be able to authenticate the user, which would be a bad user experience.
func (c *Core) Authenticate(ctx context.Context, username string, password string) (User, error) {
	usr, err := c.storer.QueryByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// Simulate a slow query to prevent timing attacks when the user does not exist.
			// These attacks are based on the fact that the query will return faster if the user does not exist.
			// We want to make sure that the query always takes the same amount of time.
			// See: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				return User{}, fmt.Errorf("rand.Int: %w", err)
			}
			time.Sleep(time.Duration(n.Int64()) * time.Millisecond)

			// Return a generic error to prevent user enumeration, but also return the underlying error so we can log it.
			return User{}, errors.Join(ErrAuthFailed, err)
		}
		return User{}, fmt.Errorf("query: username[%s]: %w", username, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		// Comparison will fail faster if the password is incorrect.
		// Add a small delay to prevent timing attacks.
		// See: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#compare-password-hashes-using-safe-functions

		n, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return User{}, fmt.Errorf("rand.Int: %w", err)
		}
		time.Sleep(time.Duration(n.Int64()) * time.Millisecond)

		return User{}, ErrAuthFailed
	}
	return usr.User, nil
}

func (c *Core) ResetPassword(ctx context.Context, newPassword Password, userID int) error {
	pwHash, err := newPassword.Hash()
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := c.storer.AddPassword(ctx, userID, pwHash); err != nil {
		return fmt.Errorf("add password: %w", err)
	}
	return nil
}

// ResetPasswordWithToken resets the password of the user if the token is valid.
func (c *Core) ResetPasswordWithToken(ctx context.Context, token string, newPassword Password) error {
	userID, err := c.verifyResetToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	usr, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("query by email: %w", err)
	}

	if err := c.ResetPassword(ctx, newPassword, usr.ID); err != nil {
		return fmt.Errorf("reset password: %w", err)
	}

	return nil
}

// RequestPasswordResetToken generates a password reset token for the user with the given email address.
func (c *Core) RequestPasswordResetToken(ctx context.Context, email mail.Address) (string, error) {
	usr, err := c.storer.QueryByEmail(ctx, email.Address)
	if err != nil {
		return "", fmt.Errorf("query by email: %w", err)
	}

	expiry := time.Now().Add(pwResetExpiresIn)
	expLocal, _ := time.LoadLocation("Asia/Taipei")
	expiry = expiry.In(expLocal)

	token, err := c.signResetToken(usr.ID, expiry)
	if err != nil {
		return "", fmt.Errorf("sign reset token: %w", err)
	}

	return token, nil
}

// signResetToken signs a reset token with the user's email address.
func (c *Core) signResetToken(userID int, expiry time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("gp-house"),
		Subject:   fmt.Sprintf("%d", userID),
		Audience:  []string{"reset-password"},
		ExpiresAt: jwt.NewNumericDate(expiry),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString(c.secret)
}

// verifyResetToken verifies the reset token and returns the user id.
func (c *Core) verifyResetToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return c.secret, nil
	},
		jwt.WithIssuer("gp-house"),
		jwt.WithAudience("reset-password"),
	)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("unexpected claims")
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid subject: %w", err)
	}

	return id, nil
}

func (c *Core) UpdateUserTypes(ctx context.Context, updates []UpdateUserType) error {
	usrs := make([]User, len(updates))

	for i, u := range updates {
		usrs[i] = User{
			ID:       u.UserID,
			UserType: u.UserType,
		}
	}

	if err := c.storer.UpdateUserTypes(ctx, usrs); err != nil {
		return fmt.Errorf("update user types: %w", err)
	}

	return nil
}

func (c *Core) QueryUsers(ctx context.Context) ([]User, error) {
	usrs, err := c.storer.QueryUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return usrs, nil
}

func (c *Core) UpdateMe(ctx context.Context, userID int, input UpdateMe) (User, error) {
	usr := User{
		ID:          userID,
		DisplayName: input.DisplayName,
	}

	if err := c.storer.UpdateMe(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update display name: %w", err)
	}

	result, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query user after update: %w", err)
	}
	return result, nil
}

func (c *Core) ConsumeAndRefreshEmailQuota(ctx context.Context, email mail.Address) error {
	err := c.storer.ConsumeAndRefreshEmailQuota(ctx, email.Address)
	if err != nil {
		if errors.Is(err, ErrorEmailQuota) {
			return ErrorEmailQuota
		}
		return fmt.Errorf("query by email: %w", err)
	}

	return nil
}
