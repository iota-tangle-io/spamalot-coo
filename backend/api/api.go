// this package defines all interfaces between the slaves and the coordinator
// its programmed in such a way, so that different means of transport can be used
package api

type CooMsg struct {
	Type    CooMsgType  `json:"type" bson:"typ"`
	Payload interface{} `json:"payload" bson:"payload"`
}

type CooMsgType byte

const (
	// single spammer
	CREATE_SP  CooMsgType = 0
	READ_SP    CooMsgType = 1
	UPDATE_SP  CooMsgType = 2
	RESTART_SP CooMsgType = 3

	// multiple spammers
	STOP_SPS    CooMsgType = 10
	RESTART_SPS CooMsgType = 11
	DELETE_SPS  CooMsgType = 12

	// errors
	ERR_HELLO_NOT_SENT CooMsgType = 20
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
	HELLO SlaveMsgType = 0
)

type SlaveHelloMsg struct {
	APIToken string `json:"api_token" bson:"api_token"`
}

type SlaveMsg struct {
	Type    SlaveMsgType `json:"type" bson:"type"`
	Payload interface{}
}
