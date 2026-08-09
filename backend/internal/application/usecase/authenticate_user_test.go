package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

type fakeUserRepository struct {
	users map[string]*domain.User
}

func (f *fakeUserRepository) FindByUsername(_ context.Context, username string) (*domain.User, error) {
	user, ok := f.users[username]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (f *fakeUserRepository) FindByID(_ context.Context, id int64) (*domain.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// fakeHasher treats a hash as valid when it is the plain value reversed,
// which keeps the test independent of any hashing library.
type fakeHasher struct{}

func (fakeHasher) Hash(plain string) (string, error) {
	return reverse(plain), nil
}

func (fakeHasher) Verify(hashed string, plain string) bool {
	return hashed == reverse(plain)
}

func reverse(value string) string {
	runes := []rune(value)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

type fakeTokenService struct {
	issued port.TokenClaims
}

func (f *fakeTokenService) Issue(claims port.TokenClaims) (string, error) {
	f.issued = claims
	return "issued-token-value", nil
}

func (f *fakeTokenService) Validate(_ string) (port.TokenClaims, error) {
	return f.issued, nil
}

func seedUser(t *testing.T, username string, plainPassword string) *fakeUserRepository {
	t.Helper()
	hashed, err := fakeHasher{}.Hash(plainPassword)
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}
	user, err := domain.NewUser(username, hashed, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("NewUser returned error: %v", err)
	}
	user.ID = 1
	return &fakeUserRepository{users: map[string]*domain.User{username: user}}
}

func TestAuthenticateUserIssuesTokenForValidCredentials(t *testing.T) {
	users := seedUser(t, "admin", "correct-horse")
	tokens := &fakeTokenService{}

	useCase := NewAuthenticateUser(users, fakeHasher{}, tokens, 15*time.Minute)
	token, err := useCase.Execute(context.Background(), "admin", "correct-horse")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if token == "" {
		t.Fatal("a successful authentication must return a token")
	}
	if tokens.issued.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role in claims, got %q", tokens.issued.Role)
	}
	if !tokens.issued.ExpiresAt.After(time.Now().UTC()) {
		t.Fatal("token expiry must be in the future")
	}
}

func TestAuthenticateUserRejectsWrongPassword(t *testing.T) {
	users := seedUser(t, "admin", "correct-horse")

	useCase := NewAuthenticateUser(users, fakeHasher{}, &fakeTokenService{}, 15*time.Minute)
	_, err := useCase.Execute(context.Background(), "admin", "wrong-password")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticateUserHidesUnknownAccount(t *testing.T) {
	users := seedUser(t, "admin", "correct-horse")

	useCase := NewAuthenticateUser(users, fakeHasher{}, &fakeTokenService{}, 15*time.Minute)
	_, err := useCase.Execute(context.Background(), "ghost", "correct-horse")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("an unknown account must report ErrInvalidCredentials, got %v", err)
	}
}
