package cloudpockets

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CloudPocket struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Account  string  `json:"account"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) GetAllCloudPockets(c echo.Context) error {
	logger := mlog.L(c)

	var cp CloudPocket
	err := c.Bind(&cp)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	stmt, err := h.db.Prepare("SELECT * FORM cloud_pockets")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "can't prepar query cloud_pockets statment:", err.Error())
	}
	rows, err := stmt.Query()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "can't query all cloud_pockets", err.Error())
	}
	var cloudPockets = []CloudPocket{}
	for rows.Next() {
		var c CloudPocket
		err := rows.Scan(&c.ID, &c.Name, &c.Category, &c.Currency, &c.Balance)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "can't query all cloud_pockets", err.Error())
		}
		c = CloudPocket{
			ID: c.ID, Name: c.Name, Category: c.Category, Currency: c.Currency, Balance: c.Balance,
		}
		cloudPockets = append(cloudPockets, c)
	}
	return echo.NewHTTPError(http.StatusOK, cloudPockets)
}
