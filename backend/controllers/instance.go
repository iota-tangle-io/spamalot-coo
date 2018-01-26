package controllers

import (
	"github.com/globalsign/mgo"
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"github.com/globalsign/mgo/bson"
	"time"
	"github.com/pkg/errors"
	"math/rand"
	"fmt"
)

const instancesColl = "instances"

type InstanceCtrl struct {
	Mongo *mgo.Session `inject:""`
	coll  *mgo.Collection
}

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (ctrl *InstanceCtrl) Init() error {
	ctrl.coll = ctrl.Mongo.DB("").C(instancesColl)
	rand.Seed(time.Now().Unix())
	ctrl.addDefaultInstance()
	return nil
}

const defaultInstanceName = "Default Instance"
func (ctrl *InstanceCtrl) addDefaultInstance() {
	var instance = &models.Instance{
		ID: bson.ObjectId("56d86ac3ce70d40078d36dc3"),
		Address: "127.0.0.1:16254",
		Name: defaultInstanceName,
		Desc: "default instance, usually running on the same host as the coordinator",
		Tags: []string{"default", "local"},
		Online: false, // probably not
	}
	if _, err := ctrl.ByName(instance.ID.Hex()); err != nil {
		if err := ctrl.Add(instance); err != nil {
			fmt.Printf("unable to add default instance, %v", err)
		}
	}
}

func (ctrl *InstanceCtrl) NewAPIToken() string {
	var apiToken string
	// TODO: there's a race condition here, but a low change of happening
	for {
		b := make([]byte, 32)
		for i := range b {
			b[i] = tokenChars[rand.Intn(len(tokenChars))]
		}
		apiToken = string(b)
		// ensure that no other instance has the same API token
		if _, err := ctrl.ByAPIToken(apiToken); err != nil {
			break
		}
	}
	return apiToken
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
	err := ctrl.coll.FindId(bson.M{"name": name}).One(&instance)
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
	instance.APIToken = ctrl.NewAPIToken()
	err := ctrl.coll.Insert(instance)
	return errors.WithStack(err)
}

func (ctrl *InstanceCtrl) Update(instance *models.Instance) error {
	err := ctrl.coll.UpdateId(instance.ID, bson.M{
		"$set": bson.M{
			"address": instance.Address,
			"name":       instance.Name,
			"desc":       instance.Desc,
			"tags":       instance.Tags,
			// online is set by the coordinator
			"updated_on": time.Now,
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
