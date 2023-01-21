//go:build unit
// +build unit

package cloudpockets

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocket(t *testing.T) {
	tests := []struct {
		name       string
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"create pocket succesfully",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(cStmt).WithArgs("pocket-name", 1).WillReturnRows(row)
				return db, err
			},
			`{"name": "pocket-name", "account_id" : 1}`,
			http.StatusCreated,
			`{"id": 1, "name": "pocket-name"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := tc.sqlFn()
			h := New(db)
			// Assertions
			assert.NoError(t, err)
			if assert.NoError(t, h.Create(c)) {
				assert.Equal(t, tc.wantStatus, rec.Code)
				assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}
		})
	}
}

func TestCreatePocket_Error(t *testing.T) {
	someErr := errors.New("some random error")
	tests := []struct {
		name    string
		sqlFn   func() (*sql.DB, error)
		reqBody string
		wantErr error
	}{
		{"create pocket failed",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(cStmt).WithArgs("pocket-name").WillReturnError(someErr)
				return db, err
			},
			`{"name": "pocket-name"}`,
			echo.NewHTTPError(http.StatusInternalServerError, "query row error"),
		},
		{"create with bad request",
			func() (*sql.DB, error) {
				return nil, nil
			},
			`ba`,
			echo.NewHTTPError(http.StatusBadRequest, "bad request body"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, _ := tc.sqlFn()
			h := New(db)

			berr := h.Create(c)
			// Assertions
			assert.Equal(t, berr, tc.wantErr)
		})
	}
}
