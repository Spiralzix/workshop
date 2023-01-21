package cloudpockets

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateCloudPoacket(c echo.Context) error {
	e := CloudPocket{}
	err := c.Bind(&e)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot bind request")
	}

	row := h.db.QueryRow("INSERT INTO cloud_pockets (Id, Name, Catagory, Currency, Balance) values ($1, $2, $3, $4, $5)  RETURNING id", e.ID, e.Name, e.Category, e.Currency, e.Balance)
	err = row.Scan(&e.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, e)
}
