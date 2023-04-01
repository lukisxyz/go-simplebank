package api

import (
	"net/http"

	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createAccountErrorResponse struct {
	Error string `json:"error"`
}

type createAccountSuccessResponse struct {
	Data db.Account `json:"data"`
}

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(c echo.Context) error {
	req := new(createAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if err := s.validate.Struct(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	account, err := s.store.CreateAccount(c.Request().Context(), arg)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&createAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&createAccountSuccessResponse{
			Data: account,
		},
	)
}
