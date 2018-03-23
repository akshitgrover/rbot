package funcs

import (
	"gopkg.in/mgo.v2/bson"
	"rbot/models"
)

func Insert(d Db, c string, e models.Event, err *error) {
	*err = nil
	*err = d.Session.DB("rbot").C(c).Insert(e)
}

func FindOne(d Db, c string, q bson.M, rf *models.Event, err *error) {
	*err = nil
	*err = d.Session.DB("rbot").C(c).Find(q).One(rf)
}

func FindOneProject(d Db, c string, q bson.M, p bson.M, rf *models.Event, err *error) {
	*err = nil
	*err = d.Session.DB("rbot").C(c).Find(q).Select(p).One(rf)
}

func Update(d Db, c string, q bson.M, uq bson.M, err *error) {
	*err = nil
	*err = d.Session.DB("rbot").C(c).Update(q, uq)
}

func FindOneAdmin(d Db, c string, q bson.M, rf *models.Admin, err *error) {
	*err = nil
	*err = d.Session.DB("rbot").C(c).Find(q).One(rf)
}
