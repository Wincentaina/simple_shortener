package redirect

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"

	resp "good_shortener/internal/lib/api/response"
	"good_shortener/internal/lib/logger/sl"
	"good_shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// TODO: понять правильно ли используется return
func New(log *slog.Logger, urlGetter URLGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)

		alias := c.Param("alias")
		if alias == "" {
			log.Info("alias is empty")

			return c.JSON(http.StatusBadRequest, resp.Error("invalid request"))
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)

				return c.JSON(http.StatusNotFound, resp.Error("not found"))
			}

			log.Error("failed to get url", sl.Err(err))
			return c.JSON(http.StatusInternalServerError, resp.Error("internal error"))
		}

		log.Info("got url", slog.String("url", resURL))

		// redirect to found url
		err = c.Redirect(http.StatusFound, resURL)
		if err != nil {
			log.Error("failed to redirect", sl.Err(err))
			return c.JSON(http.StatusInternalServerError, resp.Error("internal error"))
		}

		return nil
	}
}
