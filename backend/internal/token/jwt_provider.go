package token

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
)

// TokenProvider generates tokens for authenticated users.
type TokenProvider interface {
	Generate(username string, userID int64) (string, error)
}

// JWTProvider uses internal/utils to generate JWT tokens.
type JWTProvider struct{}

// NewJWTProvider sets the jwt secret and returns a provider.
func NewJWTProvider(secret []byte) *JWTProvider {
	utils.SetJWTKey(secret)
	return &JWTProvider{}
}

// Generate issues a signed JWT for the given user.
func (p *JWTProvider) Generate(username string, userID int64) (string, error) {
	return utils.GenerateJWT(username, userID)
}
