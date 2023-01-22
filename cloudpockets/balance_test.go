package cloudpockets

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseHandler(t *testing.T) {

	t.Run("Pass-condition", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"id", "name", "balance", "category", "currency", "account"}).
			AddRow("111", "Peter", 1000.00, "Travel", "THB", "test01")
		expected := `{"id":111,"name":"Peter","category":"Travel","currency":"THB","balance":1000,"account":"test01"}`

		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cloud-pockets/userTest", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		db, mock, err := sqlmock.New()
		cfg := config.FeatureFlag{}
		mock.ExpectQuery("SELECT (.+) FROM cloud_pockets").WithArgs().WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// cfg := config.FeatureFlag
		h := handler{cfg, db}
		c := e.NewContext(req, rec)

		// Act
		err = h.GetBalance(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Fail404, UserNotFound", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		db, mock, err := sqlmock.New()
		cfg := config.FeatureFlag{}
		mock.ExpectQuery("SELECT (.+) FROM cloud_pockets").WithArgs().WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{cfg, db}
		c := e.NewContext(req, rec)
		c.SetPath("/cloud-pockets/:id")
		c.SetParamNames("id")
		c.SetParamValues("userError")

		// Act
		err = h.GetBalance(c)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
