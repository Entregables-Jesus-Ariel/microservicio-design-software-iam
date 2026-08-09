package port

// PasswordHasher hashes and verifies plain-text passwords.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(hashed string, plain string) bool
}
