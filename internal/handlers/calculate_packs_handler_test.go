package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/GoCodingX/repartners/internal/handlers"
	"github.com/GoCodingX/repartners/internal/repository"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostPacksCalculate(t *testing.T) {
	t.Run("returns StatusInternalServerError when no packs exists", func(t *testing.T) {
		// create echo context
		c, _ := newEchoContext(&newEchoContextParams{
			method:  http.MethodPost,
			target:  "/packs/calculate",
			payload: `{}`,
		})

		// create service
		svc, repo := newServiceWithMockRepo(t)

		// mock
		repo.
			EXPECT().
			GetPacks(c.Request().Context()).
			Return(nil, nil)

		// act
		err := svc.PostPacksCalculate(c)

		// assert
		require.Error(t, err)

		var httpErr *echo.HTTPError

		require.ErrorAs(t, err, &httpErr)

		errRsp, ok := httpErr.Message.(*openapi.ErrorResponse)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, errRsp.Code)
		assert.Equal(t, "Bad Request", errRsp.Status)
	})

	t.Run("returns StatusOK", func(t *testing.T) {
		// create echo context
		c, rec := newEchoContext(&newEchoContextParams{
			method:  http.MethodPost,
			target:  "/packs/calculate",
			payload: `{"items":5}`,
		})

		// create service
		svc, repo := newServiceWithMockRepo(t)

		mockPacks := genRepoPack()

		// mock
		repo.
			EXPECT().
			GetPacks(c.Request().Context()).
			Return(mockPacks, nil)

		// act
		err := svc.PostPacksCalculate(c)

		// assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp *openapi.CalculatePacksResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)

		require.NoError(t, err)
		require.Len(t, *resp, 1)
		assert.Equal(t, int32(2), (*resp)[0].PackCount)
		assert.Equal(t, int32(4), (*resp)[0].PackSize)
	})
}

func genRepoPack() []repository.Pack {
	now := time.Now().UTC()

	return []repository.Pack{{
		ID:        uuid.New(),
		Size:      4,
		CreatedAt: now.Add(-5 * time.Minute),
		UpdatedAt: now,
	}}
}

func TestCalculateOrder(t *testing.T) {
	tests := []struct {
		orderedItemsCount int
		packs             []repository.Pack
		expectedResult    handlers.Result
	}{
		{
			orderedItemsCount: 10,
			packs: []repository.Pack{
				{Size: 5},
				{Size: 10},
				{Size: 15},
			},
			expectedResult: handlers.Result{
				TotalSent: 10,
				PackCount: 1,
				Packs:     map[int32]int32{10: 1},
			},
		},
		{
			orderedItemsCount: 15,
			packs: []repository.Pack{
				{Size: 5},
				{Size: 10},
			},
			expectedResult: handlers.Result{
				TotalSent: 15,
				PackCount: 2,
				Packs:     map[int32]int32{5: 1, 10: 1},
			},
		},
		{
			orderedItemsCount: 12,
			packs: []repository.Pack{
				{Size: 5},
				{Size: 10},
				{Size: 15},
			},
			expectedResult: handlers.Result{
				TotalSent: 15,
				PackCount: 1,
				Packs:     map[int32]int32{15: 1},
			},
		},
		{
			orderedItemsCount: 20,
			packs: []repository.Pack{
				{Size: 10},
				{Size: 20},
			},
			expectedResult: handlers.Result{
				TotalSent: 20,
				PackCount: 1,
				Packs:     map[int32]int32{20: 1},
			},
		},
		{
			orderedItemsCount: 7,
			packs: []repository.Pack{
				{Size: 3},
				{Size: 5},
			},
			expectedResult: handlers.Result{
				TotalSent: 8,
				PackCount: 2,
				Packs:     map[int32]int32{3: 1, 5: 1},
			},
		},
		{
			orderedItemsCount: 100,
			packs: []repository.Pack{
				{Size: 10},
				{Size: 25},
				{Size: 50},
			},
			expectedResult: handlers.Result{
				TotalSent: 100,
				PackCount: 2,
				Packs:     map[int32]int32{50: 2},
			},
		},
		{
			orderedItemsCount: 8,
			packs: []repository.Pack{
				{Size: 5},
			},
			expectedResult: handlers.Result{
				TotalSent: 10,
				PackCount: 2,
				Packs:     map[int32]int32{5: 2},
			},
		},
		{
			orderedItemsCount: 13,
			packs: []repository.Pack{
				{Size: 10},
				{Size: 20},
			},
			expectedResult: handlers.Result{
				TotalSent: 20,
				PackCount: 1,
				Packs:     map[int32]int32{20: 1},
			},
		},
		{
			orderedItemsCount: 16,
			packs: []repository.Pack{
				{Size: 8},
				{Size: 25},
			},
			expectedResult: handlers.Result{
				TotalSent: 16,
				PackCount: 2,
				Packs:     map[int32]int32{8: 2},
			},
		},
		{
			orderedItemsCount: 23,
			packs: []repository.Pack{
				{Size: 3},
				{Size: 7},
				{Size: 11},
				{Size: 15},
			},
			expectedResult: handlers.Result{
				TotalSent: 23,
				PackCount: 5,
				Packs:     map[int32]int32{3: 4, 11: 1},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			// act
			result := handlers.CalculateOrder(tt.orderedItemsCount, tt.packs)

			// assert
			assert.Equal(t, tt.expectedResult.TotalSent, result.TotalSent)
			assert.Equal(t, tt.expectedResult.PackCount, result.PackCount)
			assert.Equal(t, tt.expectedResult.Packs, result.Packs)
			assert.GreaterOrEqual(t, result.TotalSent, tt.orderedItemsCount,
				"TotalSent should be >= ordered items")
			assert.LessOrEqual(
				t,
				result.TotalSent, tt.orderedItemsCount+handlers.MaxOvershoot,
				"TotalSent should be within overshoot limit",
			)
		})
	}
}
