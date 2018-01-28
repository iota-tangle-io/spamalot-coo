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
	"encoding/json"
)

type Coordinator struct {
	Config       config.CoordinatorConfig
	InstanceCtrl *controllers.InstanceCtrl `inject:""`
	webEngine    *echo.Echo
	logger       log15.Logger
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
	c.webEngine.GET("/api", c.wsHandler)
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

	// expect the first read to be a hello from the slave
	slaveMsg := &SlaveMsg{}
	if err := ws.ReadJSON(slaveMsg); err != nil {
		connLogger.Warn("unable to read first slave message, closing conn.")
		return nil
	}

	if slaveMsg.Type != SLAVE_HELLO {
		connLogger.Warn("first message was not SLAVE_HELLO, closing conn.")
		return nil
	}

	helloMsg := &SlaveHelloMsg{}
	if err := json.Unmarshal(slaveMsg.Payload, helloMsg); err != nil {
		connLogger.Warn("unable to parse payload of SLAVE_HELLO message", "err", err.Error())
		return nil
	}

	// expect an API token from the slave
	if !lib.ValidAPIToken(helloMsg.APIToken) {
		if err := ws.WriteJSON(&CooMsg{Type: SLAVE_API_TOKEN_INVALID}); err != nil {
			connLogger.Warn("wasn't able to send SLAVE_API_TOKEN_INVALID", "err", err.Error())
		}
		return nil
	}

	coo.handleSlave(ws, helloMsg.APIToken);
	return nil
}

func (coo *Coordinator) handleSlave(slaveWsConn *websocket.Conn, apiToken string) {

	// check the slave's API token
	slave, err := coo.InstanceCtrl.ByAPIToken(apiToken)
	if err != nil {
		if err := slaveWsConn.WriteJSON(&CooMsg{Type: SLAVE_API_TOKEN_INVALID}); err != nil {
			coo.logger.Warn("unable to send SLAVE_API_TOKEN_INVALID to client", "err", err.Error())
		}
		return
	}

	slaveLogger := coo.logger.New("slave", slave.Name)

	// send the slave a warm welcome after validating its existence
	if err := slaveWsConn.WriteJSON(&CooMsg{Type: SLAVE_WELCOME}); err != nil {
		coo.logger.Warn("unable to send SLAVE_WELCOME to slave", "error", err.Error())
		return
	}

	for {

		msg := &SlaveMsg{}
		if err := slaveWsConn.ReadJSON(msg); err != nil {
			slaveLogger.Warn("unable to read new message", "err", err.Error())
			break
		}

		// router for messages
		switch msg.Type {
		case SLAVE_BYE:
			slaveLogger.Info("disconnected")
			break
		}

	}
}
