package controllers

import (
	"github.com/labstack/echo"
	"github.com/iota-tangle-io/spamalot-coo/backend/utilities"
	"gopkg.in/inconshreveable/log15.v2"
	"github.com/iota-tangle-io/spamalot-coo/backend/server/config"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/iota-tangle-io/spamalot-coo/backend/lib"
	"encoding/json"
	"github.com/iota-tangle-io/spamalot-coo/api"
	"crypto/md5"
	"encoding/hex"
)

type Coordinator struct {
	Config       config.CoordinatorConfig
	InstanceCtrl *InstanceCtrl `inject:""`
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
	slaveMsg := &api.SlaveMsg{}
	if err := ws.ReadJSON(slaveMsg); err != nil {
		connLogger.Warn("unable to read first slave msg, closing conn.")
		return nil
	}

	if slaveMsg.Type != api.SLAVE_HELLO {
		connLogger.Warn("first message was not SLAVE_HELLO, closing conn.")
		return nil
	}

	helloMsg := &api.SlaveHelloMsg{}
	if err := json.Unmarshal(slaveMsg.Payload, helloMsg); err != nil {
		connLogger.Warn("unable to parse payload of SLAVE_HELLO msg", "err", err.Error())
		return nil
	}

	// expect an API token from the slave
	if !lib.ValidAPIToken(helloMsg.APIToken) {
		if err := ws.WriteJSON(&api.CooMsg{Type: api.SLAVE_API_TOKEN_INVALID}); err != nil {
			connLogger.Warn("wasn't able to send SLAVE_API_TOKEN_INVALID msg", "err", err.Error())
		}
		return nil
	}

	coo.communicate(connLogger, ws, helloMsg.APIToken);
	return nil
}

func (coo *Coordinator) communicate(connLogger log15.Logger, slaveWsConn *websocket.Conn, apiToken string) {

	// check the slave's API token
	slave, err := coo.InstanceCtrl.ByAPIToken(apiToken)
	if err != nil {
		connLogger.Warn("API token invalid", "tokenSupplied", apiToken)
		if err := slaveWsConn.WriteJSON(&api.CooMsg{Type: api.SLAVE_API_TOKEN_INVALID}); err != nil {
			connLogger.Warn("unable to send SLAVE_API_TOKEN_INVALID msg", "err", err.Error())
		}
		return
	}

	slaveLogger := connLogger.New("slave", slave.Name)

	// marshal the current spammer config as a payload
	spammerConfigBytes, err := json.Marshal(slave.SpammerConfig)
	if err != nil {
		slaveLogger.Warn("can't marshal spammer config to bytes, canceling conn.", "err", err.Error())
		// abort connection with slave (ignore error)
		slaveWsConn.WriteJSON(&api.CooMsg{Type: api.COO_INTERNAL_ERROR})
		return
	}

	// send the slave a warm welcome after validating its existence
	// and provide the latest configuration defined for it
	welcomeMsg := &api.CooMsg{Type: api.SLAVE_WELCOME, Payload: spammerConfigBytes}
	if err := slaveWsConn.WriteJSON(welcomeMsg); err != nil {
		slaveLogger.Warn("unable to send SLAVE_WELCOME to slave", "error", err.Error())
		return
	}

	spammerStateMsg := coo.readSpammerStateMsg(slaveWsConn, slaveLogger)
	if spammerStateMsg == nil {
		return
	}

	if !verifyConfigHash(spammerConfigBytes, spammerStateMsg.ConfigHash) {
		slaveLogger.Warn("slave did not adjust spammer config to coo's, canceling conn.")
		return
	}

	slaveLogger.Info("slave's configuration hash is valid")
	slaveLogger.Info("starting spammer on slave...")

	coo.sendStartMsg(slaveWsConn, slaveLogger)

	spammerStateMsg = coo.readSpammerStateMsg(slaveWsConn, slaveLogger)
	if err != nil {
		slaveLogger.Warn("unable to read spammer state after sending SP_START", "err", err.Error())
		return
	}

	if !spammerStateMsg.Running {
		slaveLogger.Warn("expected spammer to be running after start msg, canceling conn.")
		return
	}

	slaveLogger.Info("spammer was started on slave")
	slaveLogger.Info("collecting metric data...")

	for {

		msg := &api.SlaveMsg{}
		if err := slaveWsConn.ReadJSON(msg); err != nil {
			slaveLogger.Warn("unable to read new msg", "err", err.Error())
			break
		}

		// router for messages
		switch msg.Type {
		case api.SLAVE_BYE:
			slaveLogger.Info("disconnected")
			break
		case api.SLAVE_SPAMMER_STATE:
			coo.printSlaveStateInfo(slaveWsConn, slaveLogger, msg.Payload)
		default:
			slaveLogger.Warn("got an unknown msg type from slave", "code", msg.Type)
		}

	}
}

func (coo *Coordinator) printSlaveStateInfo(slaveWsConn *websocket.Conn, logger log15.Logger, payload []byte) {
	state := &api.SlaveSpammerStateMsg{}
	if err := json.Unmarshal(payload, state); err != nil {
		logger.Warn("unable to parse payload of SLAVE_SPAMMER_STATE msg", "err", err.Error())
		return
	}
	logger.Info("got slave state msg", "running", state.Running)
}

func (coo *Coordinator) sendStartMsg(slaveWsConn *websocket.Conn, logger log15.Logger) {
	if err := slaveWsConn.WriteJSON(&api.CooMsg{Type: api.SP_START, Payload: []byte{}}); err != nil {
		logger.Warn("unable to send SP_START msg", "error", err.Error())
		return
	}
	logger.Info("sent SP_START msg")
}

func (coo *Coordinator) sendStopMsg(slaveWsConn *websocket.Conn, logger log15.Logger) {
	if err := slaveWsConn.WriteJSON(&api.CooMsg{Type: api.SP_STOP}); err != nil {
		logger.Warn("unable to send SP_STOP msg", "error", err.Error())
		return
	}
	logger.Info("sent SP_STOP msg")
}

func (coo *Coordinator) readSpammerStateMsg(slaveWsConn *websocket.Conn, logger log15.Logger) *api.SlaveSpammerStateMsg {
	msg := &api.SlaveMsg{}
	if err := slaveWsConn.ReadJSON(msg); err != nil {
		logger.Warn("unable to read expected spammer state msg", "err", err.Error())
		return nil
	}

	if msg.Type != api.SLAVE_SPAMMER_STATE {
		logger.Warn("expected SLAVE_SPAMMER_STATE msg from slave", "actualCode", msg.Type)
		return nil
	}

	spammerStateMsg := &api.SlaveSpammerStateMsg{}
	if err := json.Unmarshal(msg.Payload, spammerStateMsg); err != nil {
		logger.Warn("unable to parse payload of SLAVE_SPAMMER_STATE msg", "err", err.Error())
		return nil
	}
	return spammerStateMsg
}

func verifyConfigHash(should []byte, is string) bool {
	hasher := md5.New()
	hasher.Write(should)
	return hex.EncodeToString(hasher.Sum(nil)) == is
}
