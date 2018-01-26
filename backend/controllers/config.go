package controllers

import (
	"github.com/iota-tangle-io/spamalot-coo/backend/models"
	"github.com/pkg/errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const configColl = "config"
const configID = "config"

type ConfigCtrl struct {
	Mongo *mgo.Session `inject:""`
	coll  *mgo.Collection
}

func (ctrl *ConfigCtrl) Init() error {
	ctrl.coll = ctrl.Mongo.DB("").C(configColl)

	if _, err := ctrl.Config(); err != nil {
		// create new config
		c := &models.Config{ID: configID}
		if err := ctrl.coll.Insert(c); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (ctrl *ConfigCtrl) Config() (*models.Config, error) {
	c := &models.Config{}
	err := ctrl.coll.FindId(configID).One(c)
	return c, errors.WithStack(err)
}

func (ctrl *ConfigCtrl) Update(c models.Config) error {
	err := ctrl.coll.UpdateId(configID, bson.M{
		"$set": bson.M{},
	})
	return errors.WithStack(err)
}
