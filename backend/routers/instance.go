package routers

import (
	"github.com/labstack/echo"
	"github.com/iota-tangle-io/spamalot-coo/backend/controllers"
	"net/http"
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"fmt"
)

type InstanceRouter struct {
	WebEngine    *echo.Echo                `inject:""`
	InstanceCtrl *controllers.InstanceCtrl `inject:""`
}

func instanceRoute(instance *models.Instance) string {
	return fmt.Sprintf("/api/instances/id/%s", instance.ID.Hex())
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
		return c.Redirect(http.StatusSeeOther, instanceRoute(instance))
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

		return c.Redirect(http.StatusSeeOther, instanceRoute(instance))
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
