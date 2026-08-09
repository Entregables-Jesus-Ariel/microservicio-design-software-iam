package port

// PasswordHasher hides the hashing algorithm from the application layer.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(hashed string, plain string) bool
}
