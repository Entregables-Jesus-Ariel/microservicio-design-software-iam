// Package security implements the cryptographic ports declared by the
// application layer. No other layer imports a crypto library directly.
package security

import "golang.org/x/crypto/bcrypt"

// BcryptHasher hashes and verifies passwords with bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher builds a hasher. A zero or out-of-range cost falls back
// to the library default, which is never weaker than the recommendation.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash returns the bcrypt hash of a plain password.
func (h *BcryptHasher) Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify reports whether a plain password matches a stored hash.
func (h *BcryptHasher) Verify(hashed string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
