package token

import (
	"time"
)

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// CreateToken creates a new token for a specific user ID and duration
	CreateToken(userId int64, email string, duration time.Duration, tokenType TokenType) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}
