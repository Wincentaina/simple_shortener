package getUrls

import (
	"github.com/labstack/echo/v4"
	resp "good_shortener/internal/lib/api/response"
	"good_shortener/internal/middlewares/jwtMiddleware"
	"log/slog"
	"net/http"
)
import pgf "good_shortener/internal/storage/postgres"

type Response struct {
	resp.Response
	Urls []pgf.Url `json:"urls"`
}

type UrlsGetter interface {
	GetUserUrls(user_id int64) ([]pgf.Url, error)
}

type UserGetter interface {
	GetUserByUsername(username string) (pgf.User, error)
}

func New(log *slog.Logger, urlsGetter UrlsGetter, userGetter UserGetter) echo.HandlerFunc {
	const op = "handlers.url.getUrls.New"

	return func(c echo.Context) error {
		currentUsername := jwtMiddleware.Restricted(c)
		userInfo, err := userGetter.GetUserByUsername(currentUsername)
		if err != nil {
			log.Error("error while getting user by username")
			return c.JSON(http.StatusInternalServerError, resp.Error("bad request"))
		}
		userId := userInfo.Id

		var urls []pgf.Url
		urls, err = urlsGetter.GetUserUrls(userId)

		log.Info("successful getting urls")

		return c.JSON(http.StatusOK, Response{
			Response: resp.OK(),
			Urls:     urls,
		})
	}
}
