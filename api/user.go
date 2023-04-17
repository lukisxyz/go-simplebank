package api

import (
	"database/sql"
	"net/http"

	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

type createUserSuccessResponse struct {
	Data UserResponse `json:"data"`
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserErrorResponse struct {
	Error string `json:"error"`
}

func (r createUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(&r.Username, validation.Required,
			validation.Length(8, 32),
		),
		validation.Field(&r.Fullname, validation.Required,
			validation.Length(8, 32),
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(8, 32),
		),
	)
}

func (s *Server) CreateUser(c echo.Context) error {
	req := new(createUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	hashPassword, err := util.GenerateHashFromPassword(req.Password, s.passwordHashing)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.Fullname,
		Email:          req.Email,
		HashedPassword: hashPassword,
	}

	user, err := s.store.CreateUser(c.Request().Context(), arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.JSON(
					http.StatusForbidden,
					&createUserErrorResponse{
						Error: err.Error(),
					},
				)
			}
		}
		return c.JSON(
			http.StatusInternalServerError,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&createUserSuccessResponse{
			Data: generateUserResponse(user),
		},
	)
}

type loginSuccessResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

type loginRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username"`
}

func (r loginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required,
			validation.Length(5, 32),
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(8, 32),
		),
	)
}

func (s *Server) LoginUser(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	user, err := s.store.GetUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound,
				&createUserErrorResponse{
					Error: err.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	isMatch, err := util.ComparePasswordAndHashPassword(req.Password, user.HashedPassword, s.passwordHashing)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	if !isMatch {
		return c.JSON(
			http.StatusUnauthorized,
			&createUserErrorResponse{
				Error: util.ErrWrongPassword.Error(),
			},
		)
	}

	accessToken, _, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&createUserErrorResponse{
				Error: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		&loginSuccessResponse{
			User:        generateUserResponse(user),
			AccessToken: accessToken,
		},
	)
}

func generateUserResponse(u db.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		FullName: u.FullName,
		Email:    u.Email,
	}
}
