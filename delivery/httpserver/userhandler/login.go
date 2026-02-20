package userhandler

import (
	"net/http"
	"quiz-game/param"
	"quiz-game/pkg/httpmsg"

	"github.com/labstack/echo/v4"
)

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if fieldErrors, err := h.userValidator.ValidateLoginRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		//return echo.NewHTTPError(code, msg, fieldErrors)
		return c.JSON(code, echo.Map{
			"message": msg,
			"error":   fieldErrors,
		})
	}

	resp, err := h.userSvc.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
