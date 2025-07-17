package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GoCodingX/repartners/internal/repository"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPostPacks(t *testing.T) {
	t.Run("returns StatusConflict when pack size already exists", func(t *testing.T) {
		// create echo context
		c, _ := newEchoContext(&newEchoContextParams{
			method:  http.MethodPost,
			target:  "/packs",
			payload: `{}`,
		})

		// create service
		svc, repo := newServiceWithMockRepo(t)

		// mock
		repo.
			EXPECT().
			CreatePack(c.Request().Context(), gomock.Any()).
			Return(&repository.AlreadyExistsError{})

		// act
		err := svc.PostPacks(c)

		// assert
		require.Error(t, err)

		var httpErr *echo.HTTPError

		require.ErrorAs(t, err, &httpErr)

		errRsp, ok := httpErr.Message.(*openapi.ErrorResponse)

		require.True(t, ok)
		assert.Equal(t, http.StatusConflict, errRsp.Code)
		assert.Equal(t, "Conflict", errRsp.Status)
	})

	t.Run("returns StatusCreated", func(t *testing.T) {
		sizeInput := 300

		// create echo context
		c, rec := newEchoContext(&newEchoContextParams{
			method:  http.MethodPost,
			target:  "/packs",
			payload: fmt.Sprintf(`{"size":%d}`, sizeInput),
		})

		// create service
		svc, repo := newServiceWithMockRepo(t)

		// mock
		repo.
			EXPECT().
			CreatePack(c.Request().Context(), gomock.Any()).
			Return(nil)

		// act
		err := svc.PostPacks(c)

		// assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp *openapi.Pack
		err = json.Unmarshal(rec.Body.Bytes(), &resp)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, resp.Id)
		assert.Equal(t, int32(sizeInput), resp.Size)
	})
}
