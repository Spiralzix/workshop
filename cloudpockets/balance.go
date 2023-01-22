package cloudpockets

import (
	"database/sql"
	"fmt"
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
	row := h.db.QueryRow("SELECT ID, Name, Balance, Category, Currency FROM cloud_pockets WHERE id = $1", id)
	err = row.Scan(&e.ID, &e.Name, &e.Balance, &e.Category, &e.Currency)
	fmt.Println("Err", err)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("111")
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case nil:
		fmt.Println("222")
		return c.JSON(http.StatusCreated, e)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
