package api

import (
	"github.com/iota-tangle-io/spamalot-coo/backend/controllers"
	"github.com/labstack/echo"
	"github.com/iota-tangle-io/spamalot-coo/backend/utilities"
	"gopkg.in/inconshreveable/log15.v2"
	"github.com/iota-tangle-io/spamalot-coo/backend/server/config"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/iota-tangle-io/spamalot-coo/backend/lib"
)

type Coordinator struct {
	Config config.CoordinatorConfig
	InstanceCtrl *controllers.InstanceCtrl `inject:""`
	webEngine *echo.Echo
	logger log15.Logger
}

func (c *Coordinator) Run() {
	logger, err := utilities.GetLogger("app")
	if err != nil {
		panic(err)
	}
	c.logger = logger
	logger.Info("booting up coordinator...")

	c.webEngine = echo.New()
	c.webEngine.HideBanner = true
	c.webEngine.GET("/ws", c.wsHandler)
	listenAddress := c.Config.Address
	go c.webEngine.Start(listenAddress)

	logger.Info(fmt.Sprintf("coordinator ready, listening on %s", listenAddress))
}

var (
	upgrader = websocket.Upgrader{}
)


func (coo *Coordinator) wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	connLogger := coo.logger.New("address", ws.LocalAddr().String())

	connLogger.Info("new slave ws connection")

	for {
		// expect the first read to be a hello from the slave
		helloMsg := &SlaveHelloMsg{}
		if err := ws.ReadJSON(helloMsg); err != nil {
			connLogger.Warn("first message was not HELLO, closing conn.")
			break;
		}

		// expect an API token from the slave
		if !lib.ValidAPIToken(helloMsg.APIToken) {
			if err := ws.WriteJSON(&CooMsg{Type: ERR_HELLO_NOT_SENT}); err != nil {
				connLogger.Warn("wasn't able to send ERR_HELLO_NOT_SENT", "err", err.Error())
			}
			break
		}
	}

	return nil
}