package pg

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GoCodingX/repartners/internal/repository"
	"github.com/GoCodingX/repartners/pkg/db"
)

func (r *Repository) CreatePack(ctx context.Context, pack *repository.Pack) error {
	_, err := r.db.NewInsert().Model(pack).Exec(ctx)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return repository.NewAlreadyExistsError("size", strconv.Itoa(int(pack.Size)), err)
		}

		return fmt.Errorf("failed to persist pack in db: %w", err)
	}

	return nil
}

func (r *Repository) GetPacks(ctx context.Context) ([]repository.Pack, error) {
	var packs []repository.Pack

	query := r.db.NewSelect().
		Model(&packs).
		Limit(10)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return packs, nil
}

func (r *Repository) DeletePack(ctx context.Context, packID string) error {
	_, err := r.db.NewDelete().Model(&repository.Pack{}).Where("id = ?", packID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete pack from db: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePack(ctx context.Context, packID string, size int32) error {
	_, err := r.db.NewUpdate().
		Model(&repository.Pack{}).
		Set("size = ?", size).
		Where("id = ?", packID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update pack from db: %w", err)
	}

	return nil
}
