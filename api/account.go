package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/flukis/simplebank/db/sqlc"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (r createAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Owner, validation.Required),
		validation.Field(&r.Currency, validation.Required, validation.In("USD", "EUR")),
	)
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

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
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

type getAccountErrorResponse struct {
	Error string `json:"error"`
}

type getAccountSuccessResponse struct {
	Data db.Account `json:"data"`
}

func (s *Server) getAccount(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&getAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	account, err := s.store.GetAccount(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound,
				&getAccountErrorResponse{
					Error: err.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			&getAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&getAccountSuccessResponse{
			Data: account,
		},
	)
}

type fetchAccountErrorResponse struct {
	Error string `json:"error"`
}

type fetchAccountSuccessResponse struct {
	Data []db.Account `json:"data"`
	Meta Meta         `json:"meta"`
}

type fetchAccountRequest struct {
	PageID int32 `form:"page" binding:"required"`
	Limit  int32 `form:"limit" binding:"required"`
}

func (r fetchAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PageID, validation.Required, validation.Min(1)),
		validation.Field(&r.Limit, validation.Required, validation.Min(5), validation.Max(50), validation.MultipleOf(5)),
	)
}

func (s *Server) FetchAccount(c echo.Context) error {
	req := new(fetchAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&fetchAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&fetchAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	arg := db.FetchAccountsParams{
		Limit:  req.Limit,
		Offset: (req.PageID - 1) * req.Limit,
	}

	account, err := s.store.FetchAccounts(c.Request().Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound,
				&fetchAccountErrorResponse{
					Error: err.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			&fetchAccountErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if len(account) == 0 {
		return c.JSON(
			http.StatusNotFound,
			&fetchAccountErrorResponse{
				Error: "record for accounts not found",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&fetchAccountSuccessResponse{
			Data: account,
			Meta: Meta{
				Limit: req.Limit,
				Page:  req.PageID,
			},
		},
	)
}
