// this package defines all interfaces between the slaves and the coordinator
// its programmed in such a way, so that different means of transport can be used
package api

import "encoding/json"

type CooMsg struct {
	Type    CooMsgType      `json:"type" bson:"typ"`
	Payload json.RawMessage `json:"payload" bson:"payload"`
}

type CooMsgType byte

const (
	UNDEFINED CooMsgType = 0

	// single spammer
	CREATE_SP  CooMsgType = 1
	READ_SP    CooMsgType = 2
	UPDATE_SP  CooMsgType = 3
	RESTART_SP CooMsgType = 4

	// multiple spammers
	STOP_SPS    CooMsgType = 10
	RESTART_SPS CooMsgType = 11
	DELETE_SPS  CooMsgType = 12

	// errors
	SLAVE_API_TOKEN_INVALID CooMsgType = 20

	// slave
	SLAVE_WELCOME CooMsgType = 30
)

type PoWMode byte

const (
	POW_LOCAL  PoWMode = 0
	POW_REMOTE PoWMode = 0
)

type SpammerConfig struct {
	NodeAddress string  `json:"node_address" bson:"node_address"`
	SecurityLvl byte    `json:"security_lvl" bson: "security_lvl"`
	MWM         byte    `json:"mwm" bson:"mwm"`
	Depth       byte    `json:"depth" bson:"depth"`
	Tag         string  `json:"tag" bson:"tag"`
	Message     string  `json:"message" bson:"message"`
	PoWMode     PoWMode `json:"pow_mode" bson:"pow_mode"`
}

type SlaveMsgType byte

const (
	SLAVE_HELLO SlaveMsgType = 0
	SLAVE_BYE   SlaveMsgType = 1
)

type SlaveHelloMsg struct {
	APIToken string `json:"api_token" bson:"api_token"`
}

type SlaveMsg struct {
	Type    SlaveMsgType    `json:"type" bson:"type"`
	Payload json.RawMessage `json:"payload" bson:"payload"`
}
