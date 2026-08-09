package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

type fakeMovementRepository struct {
	created  []*domain.Movement
	stored   map[int64]*domain.Movement
	totals   []port.CategoryTotal
	nextID   int64
	listArgs port.MovementFilter
}

func newFakeMovementRepository() *fakeMovementRepository {
	return &fakeMovementRepository{stored: map[int64]*domain.Movement{}, nextID: 1}
}

func (f *fakeMovementRepository) Create(_ context.Context, movement *domain.Movement) (*domain.Movement, error) {
	movement.ID = f.nextID
	f.nextID++
	f.created = append(f.created, movement)
	f.stored[movement.ID] = movement
	return movement, nil
}

func (f *fakeMovementRepository) Update(_ context.Context, movement *domain.Movement) error {
	f.stored[movement.ID] = movement
	return nil
}

func (f *fakeMovementRepository) FindByID(_ context.Context, id int64) (*domain.Movement, error) {
	movement, ok := f.stored[id]
	if !ok {
		return nil, domain.ErrMovementNotFound
	}
	return movement, nil
}

func (f *fakeMovementRepository) List(_ context.Context, filter port.MovementFilter) ([]*domain.Movement, error) {
	f.listArgs = filter
	result := make([]*domain.Movement, 0, len(f.stored))
	for _, movement := range f.stored {
		result = append(result, movement)
	}
	return result, nil
}

func (f *fakeMovementRepository) TotalsByCategoryType(_ context.Context, _ domain.Period) ([]port.CategoryTotal, error) {
	return f.totals, nil
}

type fakeCategoryRepository struct {
	byID    map[int64]*domain.Category
	byName  map[string]*domain.Category
	created []*domain.Category
}

func newFakeCategoryRepository() *fakeCategoryRepository {
	return &fakeCategoryRepository{byID: map[int64]*domain.Category{}, byName: map[string]*domain.Category{}}
}

func (f *fakeCategoryRepository) Create(_ context.Context, category *domain.Category) (*domain.Category, error) {
	category.ID = int64(len(f.byID) + 1)
	f.byID[category.ID] = category
	f.byName[category.Name] = category
	f.created = append(f.created, category)
	return category, nil
}

func (f *fakeCategoryRepository) FindByID(_ context.Context, id int64) (*domain.Category, error) {
	category, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

func (f *fakeCategoryRepository) FindByName(_ context.Context, name string) (*domain.Category, error) {
	category, ok := f.byName[name]
	if !ok {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

func (f *fakeCategoryRepository) ListByType(_ context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	result := make([]*domain.Category, 0)
	for _, category := range f.byID {
		if category.Type == categoryType {
			result = append(result, category)
		}
	}
	return result, nil
}

func (f *fakeCategoryRepository) ListAll(_ context.Context) ([]*domain.Category, error) {
	result := make([]*domain.Category, 0, len(f.byID))
	for _, category := range f.byID {
		result = append(result, category)
	}
	return result, nil
}

type fakeAuditRepository struct {
	entries []*domain.MovementAudit
}

func (f *fakeAuditRepository) Append(_ context.Context, entry *domain.MovementAudit) error {
	f.entries = append(f.entries, entry)
	return nil
}

func (f *fakeAuditRepository) ListByMovement(_ context.Context, movementID int64) ([]*domain.MovementAudit, error) {
	result := make([]*domain.MovementAudit, 0)
	for _, entry := range f.entries {
		if entry.MovementID == movementID {
			result = append(result, entry)
		}
	}
	return result, nil
}

func seedCategory(t *testing.T, repository *fakeCategoryRepository, name string, categoryType domain.CategoryType) *domain.Category {
	t.Helper()
	category, err := domain.NewCategory(name, categoryType)
	if err != nil {
		t.Fatalf("NewCategory returned error: %v", err)
	}
	stored, err := repository.Create(context.Background(), category)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	return stored
}

func TestRecordMovementPersistsMovementAndAudit(t *testing.T) {
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}
	category := seedCategory(t, categories, "Sales", domain.CategoryIncome)

	useCase := NewRecordMovement(movements, categories, audits)
	stored, err := useCase.Execute(context.Background(), RecordMovementInput{
		UserID:      1,
		CategoryID:  category.ID,
		Date:        time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
		AmountCents: 25000,
		Note:        "daily sales",
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stored.ID == 0 {
		t.Fatal("stored movement must receive an identifier")
	}
	if len(audits.entries) != 1 {
		t.Fatalf("expected one audit entry, got %d", len(audits.entries))
	}
	if audits.entries[0].Action != domain.AuditCreate {
		t.Fatalf("expected create action, got %s", audits.entries[0].Action)
	}
}

func TestRecordMovementRejectsUnknownCategory(t *testing.T) {
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}

	useCase := NewRecordMovement(movements, categories, audits)
	_, err := useCase.Execute(context.Background(), RecordMovementInput{
		UserID:      1,
		CategoryID:  99,
		Date:        time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
		AmountCents: 25000,
	})
	if !errors.Is(err, domain.ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
	if len(movements.created) != 0 {
		t.Fatal("no movement may be stored when the category is unknown")
	}
}

func TestRecordMovementRejectsNonPositiveAmount(t *testing.T) {
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}
	category := seedCategory(t, categories, "Sales", domain.CategoryIncome)

	useCase := NewRecordMovement(movements, categories, audits)
	_, err := useCase.Execute(context.Background(), RecordMovementInput{
		UserID:      1,
		CategoryID:  category.ID,
		Date:        time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
		AmountCents: 0,
	})
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}
