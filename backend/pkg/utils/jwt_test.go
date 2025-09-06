package utils_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

const testSecret = "mysecret"

// TestJWTGenerationAndValidation checks that a JWT can be generated and
// then successfully validated, and that the claims match the expected input.
func TestJWTGenerationAndValidation(t *testing.T) {
	utils.SetJWTKey([]byte(testSecret))

	username := "emcd84"
	userID := int64(58174912)

	token, err := utils.GenerateJWT(username, userID)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	if claims.Username != username || claims.UserID != userID {
		t.Errorf("unexpected claims: got %+v, want username=%s userID=%d", claims, username, userID)
	}

	// Token should expire roughly 24h from issuance.
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		t.Errorf("token already expired: remaining=%v", remaining)
	}
	if remaining > 24*time.Hour+2*time.Minute {
		t.Errorf("token expires too far in the future: %v", remaining)
	}
}

// TestInvalidSignature verifies that a JWT signed with one key cannot be
// validated using a different key, and that an appropriate error is returned.
func TestInvalidSignature(t *testing.T) {
	// Generate token with one key...
	utils.SetJWTKey([]byte("correctkey"))
	token, err := utils.GenerateJWT("emcd84", 12345678)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// ...then validate with a different key
	utils.SetJWTKey([]byte("wrongkey"))
	_, err = utils.ValidateJWT(token)
	if err == nil {
		t.Error("expected error for invalid signature, got none")
	} else {
		var ve *jwt.ValidationError
		if ok := strings.Contains(err.Error(), "signature"); !ok {
			// fallback: check if it's a ValidationError and not due to expiry
			if errors.As(err, &ve) {
				if ve.Errors&jwt.ValidationErrorSignatureInvalid == 0 {
					t.Errorf("expected signature invalid error, got: %v", err)
				}
			} else {
				t.Logf("received error (not a ValidationError): %v", err)
			}
		}
	}
}

// TestExpiredToken ensures that a JWT with an expired timestamp is correctly
// rejected by the validation function.
func TestExpiredToken(t *testing.T) {
	utils.SetJWTKey([]byte("key"))

	claims := &utils.Claims{
		Username: "emcd84",
		UserID:   87654321,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // already expired
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	// Sign with the same key used in SetJWTKey for clarity
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("key"))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = utils.ValidateJWT(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got none")
	} else {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired == 0 {
				t.Errorf("expected expired token error, got: %v", err)
			}
		} else {
			t.Errorf("expected ValidationError for expired token, got: %v", err)
		}
	}
}
