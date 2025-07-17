package handlers

import (
	"math"
	"net/http"

	"github.com/GoCodingX/repartners/internal/repository"
	"github.com/GoCodingX/repartners/pkg/errors"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/labstack/echo/v4"
)

const MaxOvershoot = 1000

type Result struct {
	TotalSent int
	PackCount int
	Packs     map[int32]int32
}

func (s *PacksService) PostPacksCalculate(c echo.Context) error {
	// read the request payload
	calculatePacksPayload := new(openapi.CalculatePacksRequest)
	if err := c.Bind(calculatePacksPayload); err != nil {
		return err
	}

	// get pack sizes the repo layer
	packs, err := s.repo.GetPacks(c.Request().Context())
	if err != nil {
		return err
	}

	if len(packs) == 0 {
		return errors.NewEchoErrorResponse(http.StatusBadRequest, "No pack sizes found, add some.", nil)
	}

	// calculate the optimal order
	result := CalculateOrder(calculatePacksPayload.Items, packs)

	// prepare http response payload
	response := make(openapi.CalculatePacksResponse, len(result.Packs))
	i := 0

	for packSize, packCount := range result.Packs {
		response[i] = openapi.PackAndItems{
			PackCount: packCount,
			PackSize:  packSize,
		}

		i++
	}

	// respond
	return c.JSON(http.StatusOK, &response)
}

//nolint:cyclop
func CalculateOrder(orderedNrOfItems int, packs []repository.Pack) Result {
	maxLimit := orderedNrOfItems + MaxOvershoot

	type dpEntry struct {
		packCount int32
		packs     map[int32]int32
	}

	dp := make([]*dpEntry, maxLimit+1)
	dp[0] = &dpEntry{
		packCount: 0,
		packs:     make(map[int32]int32),
	}

	// process to ensure we find optimal solutions
	for i := 0; i <= maxLimit; i++ {
		if dp[i] == nil {
			continue
		}

		for _, p := range packs {
			next := i + int(p.Size)
			if next > maxLimit {
				continue
			}

			newPackCount := dp[i].packCount + 1

			// check if this is a better solution
			if dp[next] == nil || dp[next].packCount > newPackCount {
				// deep copy the packs map
				newPacks := make(map[int32]int32, len(dp[i].packs)+1)
				for k, v := range dp[i].packs {
					newPacks[k] = v
				}

				newPacks[p.Size]++

				dp[next] = &dpEntry{
					packCount: newPackCount,
					packs:     newPacks,
				}
			}
		}
	}

	// find the best solution
	best := Result{
		TotalSent: math.MaxInt,
		PackCount: math.MaxInt,
		Packs:     nil,
	}

	for i := orderedNrOfItems; i <= maxLimit; i++ {
		if dp[i] != nil {
			totalSent := i
			packCount := dp[i].packCount

			// prioritize: 1. minimum total sent 2. minimum pack count
			if totalSent < best.TotalSent ||
				(totalSent == best.TotalSent && int(packCount) < best.PackCount) {
				best = Result{
					TotalSent: totalSent,
					PackCount: int(packCount),
					Packs:     dp[i].packs,
				}
			}
		}
	}

	return best
}
