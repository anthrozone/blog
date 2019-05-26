package platform

import (
	"github.com/globalsign/mgo"
)

type Platform struct {
	Mongo *mgo.Session
	Key   string
}
