package channels

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type addMessageResult struct {
	Result bool `json:"result"`
}

func HandleAddMessage(ctx echo.Context) error {
	channel := ctx.Param("channel")
	message := ctx.FormValue("message")

	result, _ := AddMessage(ctx.Request().Context(), channel, message)

	return ctx.JSON(http.StatusOK, &addMessageResult{
		result,
	})
}
