package save

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	resp "good_shortener/internal/lib/api/response"
	"good_shortener/internal/lib/logger/sl"
	"good_shortener/internal/lib/random"
	"good_shortener/internal/middlewares/jwtMiddleware"
	"good_shortener/internal/storage"
	pgEf "good_shortener/internal/storage/postgres"
	"log/slog"
	"net/http"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string, userId int64) error
}

type UserGetter interface {
	GetUserByUsername(username string) (pgEf.User, error)
}

func New(log *slog.Logger, urlSaver URLSaver, userGetter UserGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get id by username
		currentUsername := jwtMiddleware.Restricted(c)
		userInfo, err := userGetter.GetUserByUsername(currentUsername)
		userId := userInfo.Id

		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			return c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		}

		log.Info("request decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(validateErr))

			return c.JSON(http.StatusBadRequest, resp.ValidationError(validateErr))
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		// TODO: add check for existing url
		err = urlSaver.SaveURL(req.Url, alias, userId)

		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url is already exists", slog.String("url", req.Url))

				return c.JSON(http.StatusInternalServerError, resp.Error("url already exists"))
			}

			log.Error("failed to add url", sl.Err(err))
			return c.JSON(http.StatusInternalServerError, resp.Error("failed to add url"))
		}

		log.Info("url added")

		return c.JSON(http.StatusCreated, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
