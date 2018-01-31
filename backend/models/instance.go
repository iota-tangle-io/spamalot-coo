package models

import (
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/iota-tangle-io/spamalot-coo/api"
)

type Instances []*Instance

// an Slave is a slave to the coordinator which handles spammers on its host system.
//	- controls a spammer local to it
//	- communicates with the coordinator for configuration
//	- supplies its health to the coordinator
type Instance struct {
	ID            bson.ObjectId             `json:"id" bson:"_id"`
	Address       string                    `json:"address" bson:"address"`
	APIToken      string                    `json:"api_token" bson:"api_token"`
	Name          string                    `json:"name" bson:"name"`
	Desc          string                    `json:"desc" bson:"desc"`
	Tags          []string                  `json:"tags" bson:"tags"`
	Online        bool                      `json:"online" bson:"online"`
	CheckAddress  bool                      `json:"check_address" bson:"check_address"`
	SpammerConfig *api.SpammerConfig        `json:"spammer_config" bson:"spammer_config"`
	LastState     *api.SlaveSpammerStateMsg `json:"last_state" bson:"last_state"`
	CreatedOn     time.Time                 `json:"created_on" bson:"created_on"`
	UpdatedOn     time.Time                 `json:"updated_on" bson:"updated_on"`
}
