//go:build integration
// +build integration

package cloudpockets

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocketIT(t *testing.T) {
	e := echo.New()

	// cfg := config.New().All()
	sql, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		t.Error(err)
	}

	couldPocket := New(sql)
	e.POST("/cloud-pockets", couldPocket.Create)

	reqBody := `{"name": "test", "account_id" : 1}`
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{"id": 1, "name": "test"}`
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
