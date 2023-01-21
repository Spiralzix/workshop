package transaction

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	From        string    `json:"from"`
	To          string    `json:"to"`
}

type ResultTblTransaction struct {
	Id     string    `json:"id"`
	RefId  time.Time `json:"refId"`
	PkId   string    `json:"pkId"`
	Date   float64   `json:"date"`
	Desc   string    `json:"desc"`
	Amount string    `json:"amount"`
	Type   string    `json:"type"`
}

type ResponseTransaction struct {
	Id           string                 `json:"id"`
	Transactions []ResultTblTransaction `json:"transaction"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) GetTransactionById(c echo.Context) error {
	logger := mlog.L(c)

	paramId := c.Param("id")

	resTb := Transaction{}

	rows, err := h.db.Query("SELECT id, title FROM xxx WHERE id=$1", paramId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	resBody := ResponseTransaction{}
	resBody.Id = paramId

	for rows.Next() {

		err := rows.Scan(&resTb.ID, &resTb.Timestamp, &resTb.Description, &resTb.Amount, &resTb.Type, &resTb.From, &resTb.To)
		if err != nil {
			logger.Error("query row error", zap.Error(err))
			return err
		}

		// resBody.Transactions = append(resBody.Transactions, ResultTblTransaction{
		// 	ID:          resTb.ID,
		// 	Timestamp:   resTb.Timestamp,
		// 	Description: resTb.Description,
		// 	Amount:      resTb.Amount,
		// 	Type:        resTb.Type,
		// 	From:        resTb.From,
		// 	To:          resTb.To,
		// })
	}

	logger.Info("create successfully")

	return c.JSON(http.StatusCreated, resBody)
}
