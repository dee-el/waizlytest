package authservice

// Authz represents the authentication service interface.
// Intended as a provider
type Authz interface {

	// ValidateToken validates the provided token and returns the user's ID.
	// It takes a token string as input and returns the user's ID as an int64 if
	// the token is valid, or an error if validation fails.
	ValidateToken(token string) (int64, error)
}
