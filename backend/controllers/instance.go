package controllers

import (
	"github.com/globalsign/mgo"
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"github.com/globalsign/mgo/bson"
	"time"
	"github.com/pkg/errors"
	"math/rand"
	"fmt"
	"github.com/iota-tangle-io/spamalot-coo/backend/lib"
	"github.com/iota-tangle-io/spamalot-coo/api"
	"sync"
	"encoding/json"
)

const instancesColl = "instances"

type InstanceCtrl struct {
	Mongo     *mgo.Session `inject:""`
	coll      *mgo.Collection
	Coo       *Coordinator `inject:""`
	muGateway sync.RWMutex
	gateways  map[string]chan interface{}
}

func (ctrl *InstanceCtrl) Init() error {
	ctrl.coll = ctrl.Mongo.DB("").C(instancesColl)
	rand.Seed(time.Now().Unix())
	ctrl.gateways = map[string]chan interface{}{}
	//ctrl.addDefaultInstance()
	return nil
}

const defaultInstanceName = "Default Instance"

func (ctrl *InstanceCtrl) addDefaultInstance() {
	var instance = &models.Instance{
		ID:            bson.ObjectId("56d86ac3ce70d40078d36dc3"),
		Address:       "127.0.0.1:16254",
		Name:          defaultInstanceName,
		Desc:          "default instance, usually running on the same host as the coordinator",
		Tags:          []string{"default", "local"},
		SpammerConfig: api.NewDefaultSpammerConfig(),
		Online:        false, // probably not
	}
	if _, err := ctrl.ByName(instance.Name); err != nil {
		if err := ctrl.Add(instance); err != nil {
			fmt.Printf("unable to add default instance, %v", err)
		}
	}
}

func (ctrl *InstanceCtrl) All() (models.Instances, error) {
	instances := models.Instances{}
	err := ctrl.coll.Find(nil).All(&instances)
	return instances, errors.WithStack(err)
}

func (ctrl *InstanceCtrl) ByID(id string) (*models.Instance, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.Wrapf(ErrInvalidObjectId, "%s is not a object id", id)
	}
	instance := &models.Instance{}
	err := ctrl.coll.FindId(bson.ObjectIdHex(id)).One(&instance)
	return instance, errors.WithStack(err)
}

func (ctrl *InstanceCtrl) ByName(name string) (*models.Instance, error) {
	instance := &models.Instance{}
	err := ctrl.coll.Find(bson.M{"name": name}).One(&instance)
	return instance, errors.WithStack(err)
}

func (ctrl *InstanceCtrl) ByAPIToken(token string) (*models.Instance, error) {
	instance := &models.Instance{}
	err := ctrl.coll.Find(bson.M{"api_token": token}).One(&instance)
	return instance, errors.WithStack(err)
}

func (ctrl *InstanceCtrl) Add(instance *models.Instance) error {
	instance.ID = bson.NewObjectId()
	instance.CreatedOn = time.Now()

	var apiToken string
	for {
		apiToken = lib.NewAPIToken()
		// ensure that no other instance has the same API token
		if _, err := ctrl.ByAPIToken(apiToken); err != nil {
			break
		}
	}

	if instance.SpammerConfig == nil {
		instance.SpammerConfig = api.NewDefaultSpammerConfig()
	}
	instance.APIToken = apiToken
	err := ctrl.coll.Insert(instance)
	return errors.WithStack(err)
}

func (ctrl *InstanceCtrl) Update(instance *models.Instance) error {
	if err := ctrl.coll.UpdateId(instance.ID, bson.M{
		"$set": bson.M{
			"address":        instance.Address,
			"name":           instance.Name,
			"desc":           instance.Desc,
			"check_address":  instance.CheckAddress,
			"tags":           instance.Tags,
			"spammer_config": instance.SpammerConfig,
			// online is set by the coordinator
			"updated_on": time.Now(),
		},
	}); err != nil {
		return errors.WithStack(err)
	}

	// TODO: only reset config on slave if it actually changed
	spammerConfigBytes, err := json.Marshal(instance.SpammerConfig)
	if err == nil {
		if err := ctrl.SendCooMsgToSlave(instance.ID.Hex(), api.SP_RESET_CONFIG, spammerConfigBytes); err != nil {
			fmt.Printf("unable to send SP_RESET_CONFIG to slave, err: %s", err.Error())
		}
	}

	return nil
}

func (ctrl *InstanceCtrl) UpdateOnlineState(id string, online bool) error {
	if !bson.IsObjectIdHex(id) {
		return errors.Wrapf(ErrInvalidObjectId, "%s is not a object id", id)
	}
	err := ctrl.coll.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"online":     online,
			"updated_on": time.Now(),
		},
	})
	return errors.WithStack(err)
}

func (ctrl *InstanceCtrl) UpdateLastState(id string, lastState *api.SlaveSpammerStateMsg) error {
	if !bson.IsObjectIdHex(id) {
		return errors.Wrapf(ErrInvalidObjectId, "%s is not a object id", id)
	}
	err := ctrl.coll.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"last_state": lastState,
			"updated_on": time.Now(),
		},
	})
	return errors.WithStack(err)
}

func (ctrl *InstanceCtrl) Delete(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.Wrapf(ErrInvalidObjectId, "%s is not a object id", id)
	}
	err := ctrl.coll.RemoveId(bson.ObjectIdHex(id))
	return errors.WithStack(err)
}

var ErrInstanceIsOffline = errors.New("can't send command as instance is offline")

func (ctrl *InstanceCtrl) SendCooMsgToSlave(id string, cooMsgType api.CooMsgType, payload []byte) error {
	ctrl.muGateway.RLock()
	defer ctrl.muGateway.RUnlock()
	gateway, ok := ctrl.gateways[id]
	if !ok {
		return ErrInstanceIsOffline
	}
	gateway <- &api.CooMsg{Type: cooMsgType, Payload: payload}
	return nil
}

// coordinator controller communication

func (ctrl *InstanceCtrl) AddGateway(id string, gateway chan interface{}) {
	ctrl.muGateway.Lock()
	defer ctrl.muGateway.Unlock()
	_, ok := ctrl.gateways[id]
	if ok {
		return
	}
	ctrl.gateways[id] = gateway
}

func (ctrl *InstanceCtrl) RemoveGateway(id string) {
	ctrl.muGateway.Lock()
	defer ctrl.muGateway.Unlock()
	_, ok := ctrl.gateways[id]
	if !ok {
		return
	}
	delete(ctrl.gateways, id)
}
