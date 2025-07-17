package handlers

import (
	"net/http"

	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/labstack/echo/v4"
)

func (s *PacksService) GetPacks(c echo.Context) error {
	// query packs via the repo layer
	packs, err := s.repo.GetPacks(c.Request().Context())
	if err != nil {
		return err
	}

	// prepare http response payload
	response := make(openapi.GetPacksResponse, len(packs))
	for i, p := range packs {
		response[i] = openapi.Pack{
			Id:   p.ID,
			Size: p.Size,
		}
	}

	// respond
	return c.JSON(http.StatusOK, &response)
}
