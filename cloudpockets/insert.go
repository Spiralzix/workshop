package cloudpockets

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt         = "INSERT INTO cloud_pockets (Id, Name, Catagory, Currency, Balance) values ($1, $2, $3, $4, $5)  RETURNING id;"
	cBalanceLimit = 10000
)

var (
	hErrBalanceLimitExceed = echo.NewHTTPError(http.StatusBadRequest,
		"create account balance exceed limitation")
)

func (h handler) CreateCloudPoacket(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	e := CloudPocket{}
	err := c.Bind(&e)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var lastInsertId string
	err = h.db.QueryRowContext(ctx, cStmt, e.ID, e.Name, e.Category, e.Currency, e.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.String("id", lastInsertId))
	e.ID = lastInsertId
	return c.JSON(http.StatusCreated, e)
}
