package helpers

import (
	"github.com/labstack/echo/v4"
)

type HachiResponseMessage struct {
	Error   bool   `json:"error" xml:"error"`
	Message string `json:"message" xml:"message"`
}

func BindAndValidate[E any](c echo.Context) (*E, error) {
	e := new(E)
	if err := c.Bind(e); err != nil {
		return nil, err
	}
	return e, nil
}
