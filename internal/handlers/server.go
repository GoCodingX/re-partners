package handlers

import (
	"github.com/GoCodingX/repartners/internal/repository"
)

type PacksService struct {
	repo repository.Repository
}

type NewPacksServiceParams struct {
	Repo repository.Repository
}

func NewPacksService(params *NewPacksServiceParams) *PacksService {
	return &PacksService{
		repo: params.Repo,
	}
}
