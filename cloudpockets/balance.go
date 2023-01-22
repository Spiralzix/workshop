package cloudpockets

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetBalance(c echo.Context) error {
	id := c.Param("id")
	e := CloudPocket{}
	err := c.Bind(&e)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot bind request")
	}
	row := h.db.QueryRow("SELECT ID, Name, Balance, Category, Currency, Account FROM cloud_pockets WHERE id = $1", id)
	err = row.Scan(&e.ID, &e.Name, &e.Balance, &e.Category, &e.Currency, &e.Account)
	switch err {
	case sql.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case nil:
		return c.JSON(http.StatusCreated, e)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
