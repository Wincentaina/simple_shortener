package registration

import (
	"errors"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	passwd "github.com/vzglad-smerti/password_hash"
	resp "good_shortener/internal/lib/api/response"
	"good_shortener/internal/lib/logger/sl"
	"good_shortener/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Username string `json:"username"`
}

type UserSaver interface {
	SaveUser(username, passwordHash string, isAdmin bool) error
}

func New(log *slog.Logger, userSaver UserSaver) echo.HandlerFunc {
	return func(c echo.Context) error {
		const op = "handlers.auth.registration.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)

		req := Request{}
		if err := c.Bind(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			return c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		}

		log.Info("request decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(validateErr))

			return c.JSON(http.StatusBadRequest, resp.ValidationError(validateErr))
		}

		hashedPassword, err := passwd.Hash(req.Password)
		if err != nil {
			log.Error("%s: %s", op, "failed to hash password")
			return c.JSON(http.StatusInternalServerError, "internal error")
		}

		err = userSaver.SaveUser(req.Username, hashedPassword, false)
		if err != nil {
			if errors.Is(err, storage.UserExists) {
				log.Info("user with that username is already exists", slog.String("username", req.Username))

				return c.JSON(http.StatusInternalServerError, resp.Error("user already exists"))
			}
		}

		log.Info("user added")

		return c.JSON(http.StatusCreated, Response{
			Response: resp.OK(),
			Username: req.Username,
		})
	}
}
