package cloudpockets

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = "INSERT INTO cloud_pockets (Id, Name, category, Currency, Balance , Account) values ($1, $2, $3, $4, $5 , $6)  RETURNING id;"
)

func (h handler) CreateCloudPockets(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	e := CloudPocket{}
	err := c.Bind(&e)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, e.ID, e.Name, e.Category, e.Currency, e.Balance, e.Account).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	e.ID = lastInsertId
	return c.JSON(http.StatusCreated, e)
}
