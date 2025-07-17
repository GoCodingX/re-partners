package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/GoCodingX/repartners/internal/repository"
	pkgerrors "github.com/GoCodingX/repartners/pkg/errors"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *PacksService) PostPacks(c echo.Context) error {
	// read the request payload
	createPackPayload := new(openapi.CreatePackRequest)
	if err := c.Bind(createPackPayload); err != nil {
		return err
	}

	// prepare data for the repo layer
	now := time.Now().UTC()
	pack := &repository.Pack{
		ID:        uuid.New(),
		Size:      createPackPayload.Size,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// store in db via the repo layer
	err := s.repo.CreatePack(c.Request().Context(), pack)
	if err != nil {
		var errAlreadyExists *repository.AlreadyExistsError
		if errors.As(err, &errAlreadyExists) {
			return pkgerrors.NewEchoErrorResponse(http.StatusConflict, errAlreadyExists.Msg, nil)
		}

		return err
	}

	// prepare http response payload
	response := &openapi.Pack{
		Id:   pack.ID,
		Size: pack.Size,
	}

	// respond
	return c.JSON(http.StatusCreated, response)
}
