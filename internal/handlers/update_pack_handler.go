package handlers

import (
	"net/http"

	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *PacksService) UpdatePack(c echo.Context, packId openapi_types.UUID) error {
	// read the request payload
	updatePackPayload := new(openapi.UpdatePackRequest)
	if err := c.Bind(updatePackPayload); err != nil {
		return err
	}

	// update the record via the repo layer
	err := s.repo.UpdatePack(c.Request().Context(), packId.String(), updatePackPayload.Size)
	if err != nil {
		return err
	}

	// respond
	return c.JSON(http.StatusNoContent, nil)
}
