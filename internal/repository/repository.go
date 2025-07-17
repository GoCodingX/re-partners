//go:generate go run go.uber.org/mock/mockgen -package=repositorytest -source=repository.go -destination=repositorytest/repository.go .

package repository

import (
	"context"
)

type Repository interface {
	CreatePack(ctx context.Context, pack *Pack) error
	GetPacks(ctx context.Context) ([]Pack, error)
	UpdatePack(ctx context.Context, packID string, size int32) error
	DeletePack(ctx context.Context, packID string) error
}
