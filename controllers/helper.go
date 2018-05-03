package funcs

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Insert(d *mgo.Session, c string, e interface{}, err *error) {
	*err = nil
	*err = d.DB("rbot").C(c).Insert(e)
}

func FindOne(d *mgo.Session, c string, q bson.M, rf interface{}, err *error) {
	*err = nil
	*err = d.DB("rbot").C(c).Find(q).One(rf)
}

func FindAll(d *mgo.Session, c string, q bson.M, rf interface{}, err *error, dbname string) {
	*err = nil
	*err = d.DB(dbname).C(c).Find(q).All(rf)
}

func FindOneProject(d *mgo.Session, c string, q bson.M, p bson.M, rf interface{}, err *error) {
	*err = nil
	*err = d.DB("rbot").C(c).Find(q).Select(p).One(rf)
}

func Update(d *mgo.Session, c string, q bson.M, uq bson.M, err *error) {
	*err = nil
	*err = d.DB("rbot").C(c).Update(q, uq)
}

func FindOneAdmin(d *mgo.Session, c string, q bson.M, rf interface{}, err *error) {
	*err = nil
	*err = d.DB("rbot").C(c).Find(q).One(rf)
}
