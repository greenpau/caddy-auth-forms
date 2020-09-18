package local

import (
	"fmt"
	"github.com/greenpau/caddy-auth-jwt"
	"github.com/greenpau/go-identity"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"strings"
	"sync"
	"time"
)

var globalAuthenticator *Authenticator

func init() {
	globalAuthenticator = NewAuthenticator()
	return
}

// Backend represents authentication provider with local backend.
type Backend struct {
	Method        string                   `json:"method,omitempty"`
	Realm         string                   `json:"realm,omitempty"`
	Path          string                   `json:"path,omitempty"`
	TokenProvider *jwt.TokenProviderConfig `json:"-"`
	Authenticator *Authenticator           `json:"-"`
	logger        *zap.Logger
}

// NewDatabaseBackend return an instance of authentication provider
// with local backend.
func NewDatabaseBackend() *Backend {
	b := &Backend{
		Method:        "local",
		TokenProvider: jwt.NewTokenProviderConfig(),
		Authenticator: globalAuthenticator,
	}
	return b
}

// Authenticator represents database connector.
type Authenticator struct {
	db     *identity.Database
	mux    sync.Mutex
	path   string
	logger *zap.Logger
}

// NewAuthenticator returns an instance of Authenticator.
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		db: identity.NewDatabase(),
	}
}

// SetPath sets database path.
func (sa *Authenticator) SetPath(s string) {
	sa.path = s
	return
}

// CreateUser creates a user in a database
func (sa *Authenticator) CreateUser(userName, userPwd, userEmail string, userClaims map[string]interface{}) error {
	user := identity.NewUser(userName)
	if err := user.AddPassword(userPwd); err != nil {
		return fmt.Errorf("failed adding password for username %s: %s", userName, err)
	}
	if err := user.AddEmailAddress(userEmail); err != nil {
		return fmt.Errorf("failed adding email address for username %s: %s", userName, err)
	}
	if userClaims != nil {
		for k, v := range userClaims {
			if k != "roles" {
				continue
			}
			for _, role := range strings.Split(v.(string), " ") {
				if err := user.AddRole(role); err != nil {
					return fmt.Errorf("failed adding role %s for username %s: %s", role, userName, err)
				}
			}
		}
	}
	if err := sa.db.AddUser(user); err != nil {
		return fmt.Errorf("failed adding user %v to user database: %s", user, err)
	}
	if err := sa.db.SaveToFile(sa.path); err != nil {
		return fmt.Errorf("failed adding user %v, error saving database at %s: %s", user, sa.path, err)
	}

	sa.logger.Info(
		"created new user",
		zap.String("user_id", user.ID),
		zap.String("user_name", userName),
		zap.String("user_email", userEmail),
		zap.Any("user_claims", userClaims),
	)
	return nil
}

// Configure check database connectivity and required tables.
func (sa *Authenticator) Configure() error {
	sa.mux.Lock()
	defer sa.mux.Unlock()
	sa.logger.Info("local backend configuration", zap.String("db_path", sa.path))
	fileInfo, err := os.Stat(sa.path)
	if err != nil {
		if os.IsNotExist(err) {
			sa.logger.Info("local database file does not exists, creating it", zap.String("db_path", sa.path))
			if err := sa.db.SaveToFile(sa.path); err != nil {
				return fmt.Errorf("failed to create local database file at %s: %s", sa.path, err)
			}
		} else {
			return fmt.Errorf("failed obtaining information about local database file at %s: %s", sa.path, err)
		}
	}

	if fileInfo.IsDir() {
		sa.logger.Error("local database file path points to a directory", zap.String("db_path", sa.path))
		return fmt.Errorf("local database file path points to a directory")
	}
	if err := sa.db.LoadFromFile(sa.path); err != nil {
		return fmt.Errorf("failed loading local database at %s: %s", sa.path, err)
	}
	return nil
}

// AuthenticateUser checks the database for the presence of a username/email
// and password and returns user claims.
func (sa *Authenticator) AuthenticateUser(userInput, password string) (*jwt.UserClaims, int, error) {
	var user *identity.User
	var err error
	sa.mux.Lock()
	defer sa.mux.Unlock()
	if strings.Contains(userInput, "@") {
		user, err = sa.db.GetUserByEmailAddress(userInput)
	} else {
		user, err = sa.db.GetUserByUsername(userInput)
	}
	if err != nil {
		return nil, 401, fmt.Errorf("user identity not found")
	}
	if user == nil {
		return nil, 500, fmt.Errorf("user identity is nil")
	}

	userMap, authenticated, err := sa.db.AuthenticateUser(user.Username, password)
	if err != nil {
		return nil, 401, err
	}
	if !authenticated {
		return nil, 401, fmt.Errorf("authentication failed")
	}
	if userMap == nil {
		return nil, 500, fmt.Errorf("user claims is nil")
	}

	claims, err := jwt.NewUserClaimsFromMap(userMap)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to parse user claims: %s", err)
	}
	if claims.Subject == "" {
		claims.Subject = user.Username
	}

	guestRoles := map[string]bool{
		"guest":     false,
		"anonymous": false,
	}

	for _, role := range claims.Roles {
		if _, exists := guestRoles[role]; exists {
			guestRoles[role] = true
		}
	}
	for role, exists := range guestRoles {
		if !exists {
			claims.Roles = append(claims.Roles, role)
		}
	}

	return claims, 200, nil
}

// ConfigureAuthenticator configures backend for .
func (b *Backend) ConfigureAuthenticator() error {
	if b.Authenticator == nil {
		b.Authenticator = NewAuthenticator()
	}
	b.Authenticator.SetPath(b.Path)
	b.Authenticator.logger = b.logger
	if err := b.Authenticator.Configure(); err != nil {
		return err
	}

	if len(b.Authenticator.db.Users) == 0 {
		userName := "webadmin"
		userPwd := uuid.NewV4().String()
		if len(userName) < 8 || len(userPwd) < 8 {
			return fmt.Errorf("failed to create default superadmin user")
		}
		userClaims := make(map[string]interface{})
		userClaims["roles"] = "superadmin"
		userEmail := userName + "@localdomain.local"
		if err := b.Authenticator.CreateUser(userName, userPwd, userEmail, userClaims); err != nil {
			b.logger.Error("failed to create default superadmin user for the database",
				zap.String("error", err.Error()))
			return err
		}
		b.logger.Info("created default superadmin user for the database",
			zap.String("user_name", userName),
			zap.String("user_secret", userPwd),
		)
	}
	return nil
}

// ValidateConfig checks whether Backend has mandatory configuration.
func (b *Backend) ValidateConfig() error {
	if b.Path == "" {
		return fmt.Errorf("path is empty")
	}
	return nil
}

// Authenticate performs authentication.
func (b *Backend) Authenticate(reqID string, kv map[string]string) (*jwt.UserClaims, int, error) {
	if kv == nil {
		return nil, 400, fmt.Errorf("No input to authenticate")
	}
	if _, exists := kv["username"]; !exists {
		return nil, 400, fmt.Errorf("No username found")
	}
	if _, exists := kv["password"]; !exists {
		return nil, 401, fmt.Errorf("No password found")
	}
	if b.Authenticator == nil {
		return nil, 500, fmt.Errorf("local backend is nil")
	}
	claims, statusCode, err := b.Authenticator.AuthenticateUser(kv["username"], kv["password"])
	if statusCode == 200 {
		claims.Origin = b.TokenProvider.TokenOrigin
		claims.ExpiresAt = time.Now().Add(time.Duration(b.TokenProvider.TokenLifetime) * time.Second).Unix()
		return claims, statusCode, nil
	}
	return nil, statusCode, err
}

// Validate checks whether Backend is functional.
func (b *Backend) Validate() error {
	if err := b.ValidateConfig(); err != nil {
		return err
	}
	if b.logger == nil {
		return fmt.Errorf("backend logger is nil")
	}

	b.logger.Info(
		"validating local backend",
		zap.String("db_path", b.Path),
	)

	if b.Authenticator == nil {
		return fmt.Errorf("local authenticator is nil")
	}

	return nil
}

// GetRealm return authentication realm.
func (b *Backend) GetRealm() string {
	return b.Realm
}

// ConfigureTokenProvider configures TokenProvider.
func (b *Backend) ConfigureTokenProvider(upstream *jwt.TokenProviderConfig) error {
	if upstream == nil {
		return fmt.Errorf("upstream token provider is nil")
	}
	if b.TokenProvider == nil {
		b.TokenProvider = jwt.NewTokenProviderConfig()
	}
	if b.TokenProvider.TokenName == "" {
		b.TokenProvider.TokenName = upstream.TokenName
	}
	if b.TokenProvider.TokenSecret == "" {
		b.TokenProvider.TokenSecret = upstream.TokenSecret
	}
	if b.TokenProvider.TokenIssuer == "" {
		b.TokenProvider.TokenIssuer = upstream.TokenIssuer
	}
	if b.TokenProvider.TokenOrigin == "" {
		b.TokenProvider.TokenOrigin = upstream.TokenOrigin
	}
	if b.TokenProvider.TokenLifetime == 0 {
		b.TokenProvider.TokenLifetime = upstream.TokenLifetime
	}
	return nil
}

// ConfigureLogger configures backend with the same logger as its user.
func (b *Backend) ConfigureLogger(logger *zap.Logger) error {
	if logger == nil {
		return fmt.Errorf("upstream logger is nil")
	}
	b.logger = logger
	return nil
}
