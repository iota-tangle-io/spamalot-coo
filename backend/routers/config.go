package routers

import (
	"github.com/iota-tangle-io/spamalot-coo/backend/controllers"
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"github.com/labstack/echo"
	"net/http"
)

type ConfigRouter struct {
	WebEngine *echo.Echo              `inject:""`
	ConfCtrl  *controllers.ConfigCtrl `inject:""`
}

func (cr *ConfigRouter) Init() {

	group := cr.WebEngine.Group("/api/config")

	group.GET("", func(c echo.Context) error {
		config, err := cr.ConfCtrl.Config()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, config)
	})

	group.PUT("", func(c echo.Context) error {
		con := &models.Config{}
		if err := c.Bind(con); err != nil {
			return err
		}

		if err := cr.ConfCtrl.Update(*con); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, "/api/config")
	})

}
