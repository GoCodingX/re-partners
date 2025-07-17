package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *PacksService) DeletePack(c echo.Context, packId openapi_types.UUID) error {
	// delete pack from db via the repo layer
	err := s.repo.DeletePack(c.Request().Context(), packId.String())
	if err != nil {
		return err
	}

	// respond
	return c.JSON(http.StatusNoContent, nil)
}
