// this package defines all interfaces between the slaves and the coordinator
// its programmed in such a way, so that different means of transport can be used
package api

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func NewCooMsg(kind CooMsgType, payload interface{}) (*CooMsg, error) {
	msg := &CooMsg{Type: kind}
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		msg.Payload = payloadBytes
	}
	return msg, nil
}

type CooMsg struct {
	Type    CooMsgType      `json:"type" bson:"type"`
	Payload json.RawMessage `json:"payload,omitempty" bson:"payload,omitempty"`
}

type CooMsgType byte

const (
	UNDEFINED CooMsgType = 0

	// single spammer

	// create the spammer with the payload configuration
	SP_START CooMsgType = 1
	// stop the spammer
	SP_STOP CooMsgType = 2
	// restart the spammer
	SP_RESTART CooMsgType = 3
	// pull metrics
	SP_METRICS CooMsgType = 4
	// reload config and restart spammer
	SP_RESET_CONFIG CooMsgType = 5

	// errors
	SLAVE_API_TOKEN_INVALID CooMsgType = 20

	// slave
	SLAVE_WELCOME CooMsgType = 30

	// misc
	COO_INTERNAL_ERROR CooMsgType = 40
)

type PoWMode byte

const (
	POW_LOCAL  PoWMode = 0
	POW_REMOTE PoWMode = 1
)

type SpammerConfig struct {
	NodeAddress     string  `json:"node_address" bson:"node_address"`
	SecurityLvl     byte    `json:"security_lvl" bson: "security_lvl"`
	MWM             byte    `json:"mwm" bson:"mwm"`
	Depth           byte    `json:"depth" bson:"depth"`
	Tag             string  `json:"tag" bson:"tag"`
	Message         string  `json:"message" bson:"message"`
	DestAddress     string  `json:"dest_address" bson:"dest_address"`
	PoWMode         PoWMode `json:"pow_mode" bson:"pow_mode"`
	FilterTrunk     bool    `json:"filter_trunk" bson:"filter_trunk"`
	FilterBranch    bool    `json:"filter_branch" bson:"filter_branch"`
	FilterMilestone bool    `json:"filter_milestone" bson:"filter_milestone"`
}

type SlaveSpammerStateMsg struct {
	ConfigHash string `json:"config_hash" bson:"config_hash"`
	Running    bool   `json:"running" bson:"running"`
}

const NirvanaAddress = "999999999999999999999999999999999999999999999999999999999999999999999999999999999"
const DefaultMessage = "GOSPAMMER9SPAMALOT"
const DefaultTag = "999SPAMALOT"

func NewDefaultSpammerConfig() *SpammerConfig {
	return &SpammerConfig{
		NodeAddress: "http://127.0.0.1:14265",
		SecurityLvl: 3, MWM: 14, PoWMode: POW_LOCAL,
		FilterTrunk: true, FilterBranch: true, FilterMilestone: true,
		DestAddress: NirvanaAddress,
		Depth:       3, Tag: DefaultTag, Message: DefaultMessage,
	}
}

func NewSlaveMsg(kind SlaveMsgType, payload interface{}) (*SlaveMsg, error) {
	msg := &SlaveMsg{Type: kind}
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		msg.Payload = payloadBytes
	}
	return msg, nil
}

type SlaveMsgType byte

const (
	// greeting from the slave when he connects
	SLAVE_HELLO SlaveMsgType = 0

	// closing message when the slave is shutdown to the coordinator
	SLAVE_BYE SlaveMsgType = 1

	// the slave's current configuration (such as name, description, tags)
	SLAVE_CONFIG SlaveMsgType = 2

	// a message containing the current state of the spammer
	SLAVE_SPAMMER_STATE SlaveMsgType = 3

	SLAVE_INTERNAL_ERROR SlaveMsgType = 40
)

type SlaveHelloMsg struct {
	APIToken string `json:"api_token" bson:"api_token"`
}

type SlaveMsg struct {
	Type    SlaveMsgType    `json:"type" bson:"type"`
	Payload json.RawMessage `json:"payload,omitempty" bson:"payload,omitempty"`
}
