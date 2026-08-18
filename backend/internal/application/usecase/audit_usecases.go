package usecase

import (
	"context"

	"iam/internal/application/port"
	"iam/internal/domain"
)

type AuditUsecases struct {
	audit port.AuditRepository
}

func NewAuditUsecases(audit port.AuditRepository) *AuditUsecases {
	return &AuditUsecases{audit: audit}
}

// ListAuditLogins fetches a paginated list of login attempts.
func (uc *AuditUsecases) ListAuditLogins(ctx context.Context, page, limit int) (domain.PaginatedLogins, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	logs, total, err := uc.audit.ListLoginAttempts(ctx, limit, offset)
	if err != nil {
		return domain.PaginatedLogins{}, err
	}

	return domain.PaginatedLogins{
		Items: logs,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
