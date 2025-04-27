package login

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	resp "good_shortener/internal/lib/api/response"
	"good_shortener/internal/lib/logger/sl"
	pgEf "good_shortener/internal/storage/postgres"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type Response struct {
	resp.Response
	Id       int64
	Username string
	IsAdmin  bool
	Token    string `json:"token"`
}

type UserGetter interface {
	GetUserByUsername(username string) (pgEf.User, error)
}

func New(log *slog.Logger, userGetter UserGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		const op = "handlers.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)

		req := Request{}
		if err := c.Bind(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			return c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		}

		// checking user and passwd
		usr, err := userGetter.GetUserByUsername(req.Username)

		// Set custom claims
		claims := &jwtCustomClaims{
			usr.Username,
			true,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		r := Response{
			Username: usr.Username,
			Id:       usr.Id,
			Token:    t,
			IsAdmin:  usr.IsAdmin,
		}

		return c.JSON(http.StatusOK, r)
	}
}
