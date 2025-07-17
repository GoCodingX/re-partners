package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePack(t *testing.T) {
	packID := uuid.New()

	// create echo context
	c, rec := newEchoContext(&newEchoContextParams{
		method: http.MethodDelete,
		target: fmt.Sprintf("/packs/%s", packID.String()),
	})

	// create service
	svc, repo := newServiceWithMockRepo(t)

	// mock
	repo.
		EXPECT().
		DeletePack(c.Request().Context(), packID.String()).
		Return(nil)

	// act
	err := svc.DeletePack(c, packID)

	// assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
