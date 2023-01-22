//go:build unit

package cloudpockets

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudPockets(t *testing.T) {
	cfg := config.New().All()
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockSql := "SELECT * FROM cloud_pockets"
	mockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance", "account"}).
		AddRow(1, "shopping", "shopping", "THB", 100.0, "00012910099999").
		AddRow(2, "shopping", "shopping", "THB", 200.0, "00022910099847")
	mockDB, mock, err := sqlmock.New()

	mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WillReturnRows(mockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	h := New(cfg.FeatureFlag, mockDB)
	err = h.GetAllCloudPockets(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"id":1,"name":"shopping","category":"shopping","currency":"THB","balance":100,"account":"00012910099999"},{"id":2,"name":"shopping","category":"shopping","currency":"THB","balance":200,"account":"00022910099847"}]`, strings.TrimSpace(rec.Body.String()))
	}
}
