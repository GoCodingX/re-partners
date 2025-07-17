package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePack(t *testing.T) {
	t.Run("returns StatusNoContent", func(t *testing.T) {
		sizeInput := 300
		packID := uuid.New()

		// create echo context
		c, rec := newEchoContext(&newEchoContextParams{
			method:  http.MethodPatch,
			target:  fmt.Sprintf("/packs/%s", packID.String()),
			payload: fmt.Sprintf(`{"size":%d}`, sizeInput),
		})

		// create service
		svc, repo := newServiceWithMockRepo(t)

		// mock
		repo.
			EXPECT().
			UpdatePack(c.Request().Context(), packID.String(), int32(sizeInput)).
			Return(nil)

		// act
		err := svc.UpdatePack(c, packID)

		// assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
