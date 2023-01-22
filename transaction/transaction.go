package transaction

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kkgo-software-engineering/workshop/cloudpockets"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type resultTblTx struct {
	ID     int       `json:"id"`
	RefId  string    `json:"refId"`
	PkId   int       `json:"pkId"`
	TxDate time.Time `json:"txDate"`
	TxDesc string    `json:"txdesc"`
	Amount float64   `json:"amount"`
	Type   string    `json:"type"`
}

type requestTransfer struct {
	PkSrc       int     `json:"pkSrc"`
	PkDest      int     `json:"pkDest"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type responseTransfer struct {
	TransactionId string                   `json:"transactionId"`
	PkSrc         cloudpockets.CloudPocket `json:"pkSrc"`
	PkDest        cloudpockets.CloudPocket `json:"pkDest"`
	Status        string                   `json:"status"`
}

type responseTx struct {
	ID           string        `json:"id"`
	Transactions []resultTblTx `json:"transaction"`
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

	resTb := resultTblTx{}

	rows, err := h.db.Query("SELECT id, title FROM xxx WHERE id=$1", paramId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	resBody := responseTx{}
	resBody.ID = paramId

	for rows.Next() {

		err := rows.Scan(&resTb)
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

var hErrBalanceIsNotEnough = echo.NewHTTPError(http.StatusBadRequest,
	"Your pocket balance is not enough")

func (h handler) Transfer(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var reqTf requestTransfer
	err := c.Bind(&reqTf)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var cp cloudpockets.CloudPocket
	cStmt := `select account , category , balance , currency,id,name from cloud_pockets where id=$1`
	err = h.db.QueryRowContext(ctx, cStmt, reqTf.PkSrc).Scan(&cp.Account, &cp.Category, &cp.Balance, &cp.Currency, &cp.ID, &cp.Name)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	if reqTf.Amount > cp.Balance {
		logger.Error("Cannot create transaction on transfering", zap.Error(hErrBalanceIsNotEnough))
		return hErrBalanceIsNotEnough
	}

	tx, err := h.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	transactionId := uuid.NewString()
	now := time.Now()
	var from cloudpockets.CloudPocket
	var to cloudpockets.CloudPocket

	_, execErr := tx.Exec(`INSERT INTO transactions (refid, pkid, txdate,txdesc, amount, type) VALUES ($1 ,$2 ,$3 ,$4 ,$5 ,$6)`, transactionId, reqTf.PkSrc, now, reqTf.Description, reqTf.Amount, "credit")
	if execErr != nil {
		_ = tx.Rollback()
		logger.Error("query row error1", zap.Error(execErr))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	execErr = tx.QueryRow(`UPDATE cloud_pockets SET balance = balance - $1 WHERE id = $2 RETURNING account , category , balance , currency,id,name`, reqTf.Amount, reqTf.PkSrc).Scan(&from.Account, &from.Category, &from.Balance, &from.Currency, &from.ID, &from.Name)
	if execErr != nil {
		_ = tx.Rollback()
		logger.Error("query row error2", zap.Error(execErr))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	_, execErr = tx.Exec(`INSERT INTO transactions (refid,pkid,txdate,txdesc,amount,type) VALUES($1,$2,$3,$4,$5,$6)`, transactionId, reqTf.PkDest, now, reqTf.Description, reqTf.Amount, "debit")
	if execErr != nil {
		_ = tx.Rollback()
		logger.Error("query row error3", zap.Error(execErr))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	execErr = tx.QueryRow(`UPDATE cloud_pockets SET balance = balance + $1 WHERE id = $2 RETURNING account , category , balance , currency,id,name`, reqTf.Amount, reqTf.PkDest).Scan(&to.Account, &to.Category, &to.Balance, &to.Currency, &to.ID, &to.Name)
	if execErr != nil {
		_ = tx.Rollback()
		logger.Error("query row error4", zap.Error(execErr))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("query row error5", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resBody := responseTransfer{}
	resBody.TransactionId = transactionId
	resBody.PkSrc = from
	resBody.PkDest = to
	resBody.Status = "Success"

	return c.JSON(http.StatusOK, resBody)
}
