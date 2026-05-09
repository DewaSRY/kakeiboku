package token

import (
	"time"
)

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username int64, email string, duration time.Duration, tokenType TokenType) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}
