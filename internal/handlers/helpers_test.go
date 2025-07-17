package handlers_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoCodingX/repartners/internal/handlers"
	"github.com/GoCodingX/repartners/internal/repository/repositorytest"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

type newEchoContextParams struct {
	method  string
	target  string
	payload string
}

// newEchoContext returns an echo.Context.
func newEchoContext(params *newEchoContextParams) (echo.Context, *httptest.ResponseRecorder) {
	var (
		target = "/"
		method string
		body   io.Reader
	)

	if params != nil {
		if params.method != "" {
			method = params.method
		}

		if params.target != "" {
			target = params.target
		}

		if params.payload != "" {
			body = strings.NewReader(params.payload)
		}
	}

	e := echo.New()
	req := httptest.NewRequest(method, target, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	return e.NewContext(req, rec), rec
}

func newServiceWithMockRepo(t *testing.T) (*handlers.PacksService, *repositorytest.MockRepository) {
	ctrl := gomock.NewController(t)
	repo := repositorytest.NewMockRepository(ctrl)
	svc := handlers.NewPacksService(&handlers.NewPacksServiceParams{
		Repo: repo,
	})

	return svc, repo
}
