package userhandler

import (
	"net/http"
	"quiz-game/param"
	"quiz-game/pkg/constant"
	"quiz-game/pkg/httpmsg"
	"quiz-game/service/authservice"

	"github.com/labstack/echo/v4"
)

func getClaims(c echo.Context) *authservice.Claims {
	return c.Get(constant.AuthMiddlewareContextKey).(*authservice.Claims)
}

func (h Handler) userProfile(c echo.Context) error {
	claims := getClaims(c)
	resp, err := h.userSvc.Profile(param.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
