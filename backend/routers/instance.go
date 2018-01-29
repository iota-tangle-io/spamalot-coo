package routers

import (
	"github.com/labstack/echo"
	"github.com/iota-tangle-io/spamalot-coo/backend/controllers"
	"net/http"
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"fmt"
	"github.com/iota-tangle-io/spamalot-coo/api"
	"encoding/json"
)

type InstanceRouter struct {
	WebEngine    *echo.Echo                `inject:""`
	InstanceCtrl *controllers.InstanceCtrl `inject:""`
}

func instanceRoute(id string) string {
	return fmt.Sprintf("/api/instances/id/%s", id)
}

func (router *InstanceRouter) Init() {

	group := router.WebEngine.Group("/api/instances")

	group.GET("", func(c echo.Context) error {
		instances, err := router.InstanceCtrl.All()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, instances)
	})

	group.GET("/id/:id", func(c echo.Context) error {
		id := c.Param("id")
		instance, err := router.InstanceCtrl.ByID(id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, instance)
	})

	group.GET("/id/:id/start", func(c echo.Context) error {
		id := c.Param("id")
		if err := router.InstanceCtrl.SendCooMsgToSlave(id, api.SP_START, nil); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, instanceRoute(id))
	})

	group.GET("/id/:id/stop", func(c echo.Context) error {
		id := c.Param("id")
		if err := router.InstanceCtrl.SendCooMsgToSlave(id, api.SP_STOP, nil); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, instanceRoute(id))
	})

	group.GET("/id/:id/restart", func(c echo.Context) error {
		id := c.Param("id")
		if err := router.InstanceCtrl.SendCooMsgToSlave(id, api.SP_RESTART, nil); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, instanceRoute(id))
	})

	group.POST("/id/:id/reset_config", func(c echo.Context) error {
		id := c.Param("id")

		spammerConfig := &api.SpammerConfig{}
		if err := c.Bind(spammerConfig); err != nil {
			return err
		}

		spammerConfigBytes, err := json.Marshal(spammerConfig)
		if err != nil {
			return err
		}

		if err := router.InstanceCtrl.SendCooMsgToSlave(id, api.SP_RESET_CONFIG, spammerConfigBytes); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, instanceRoute(id))
	})

	group.GET("/token/:token", func(c echo.Context) error {
		token := c.Param("token")
		instance, err := router.InstanceCtrl.ByAPIToken(token)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, instance)
	})

	group.POST("/id/:id", func(c echo.Context) error {
		instance := &models.Instance{}
		if err := c.Bind(instance); err != nil {
			return err
		}

		err := router.InstanceCtrl.Add(instance)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, instanceRoute(instance.ID.Hex()))
	})

	group.PUT("/id/:id", func(c echo.Context) error {
		instance := &models.Instance{}
		if err := c.Bind(instance); err != nil {
			return err
		}

		err := router.InstanceCtrl.Update(instance)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, instanceRoute(instance.ID.Hex()))
	})

	group.DELETE("/id/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := router.InstanceCtrl.Delete(id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, SimpleMsg{"ok"})
	})

}
