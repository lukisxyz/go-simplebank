package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/flukis/simplebank/db/sqlc"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type createTransferErrorResponse struct {
	Error string `json:"error"`
}

type createTransferSuccessResponse struct {
	Data db.TransferTxResult `json:"data"`
}

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR IDR"`
	Amount        int64  `json:"amount" binding:"requied,gt=0"`
}

func (r createTransferRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Currency, validation.Required, validation.In("USD", "EUR", "IDR")),
		validation.Field(&r.FromAccountID, validation.Required, validation.Min(1)),
		validation.Field(&r.ToAccountID, validation.Required, validation.Min(1)),
		validation.Field(&r.Amount, validation.Required, validation.Min(0)),
	)
}

func (s *Server) CreateTransfer(c echo.Context) error {
	req := new(createTransferRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createTransferErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createTransferErrorResponse{
				Error: err.Error(),
			},
		)
	}

	s.validAccount(c, req.FromAccountID, req.Currency)
	s.validAccount(c, req.ToAccountID, req.Currency)

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(c.Request().Context(), arg)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&createTransferErrorResponse{
				Error: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&createTransferSuccessResponse{
			Data: transfer,
		},
	)
}

func (s *Server) validAccount(c echo.Context, accountId int64, currency string) bool {
	account, err := s.store.GetAccount(c.Request().Context(), accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(
				http.StatusNotFound,
				&createTransferErrorResponse{
					Error: err.Error(),
				},
			)
			return false
		}
		c.JSON(
			http.StatusInternalServerError,
			&createTransferErrorResponse{
				Error: err.Error(),
			},
		)
		return false
	}

	if account.Currency != currency {
		c.JSON(
			http.StatusInternalServerError,
			&createTransferErrorResponse{
				Error: fmt.Sprintf("account with id %d have mismatch currency, expected %s got %s", accountId, account.Currency, currency),
			},
		)
		return false
	}

	return true
}
