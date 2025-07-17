package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Pack struct {
	bun.BaseModel `bun:"table:packs,alias:p"`
	ID            uuid.UUID `bun:",pk,type:uuid,notnull"`
	Size          int32     `bun:",notnull"`
	CreatedAt     time.Time `bun:",notnull"`
	UpdatedAt     time.Time `bun:",notnull"`
}
