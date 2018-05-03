package funcs

// import (
// 	"encoding/json"
// 	"gopkg.in/mgo.v2/bson"
// 	"log"
// 	"net/http"
// 	"rbot/models"
// )

// type CustmErr struct {
// 	Value string
// }

// func (e CustmErr) Error() string {
// 	return e.Value
// }

// func InsertApi(w http.ResponseWriter, req *http.Request) {
// 	var err error
// 	flagch := make(chan int)
// 	log.Println(req.URL.Path)
// 	if req.Method == "GET" {
// 		handlerError(w, CustmErr{Value: "Invalid Method"}, 405)
// 		return
// 	}
// 	err = req.ParseForm()
// 	if err != nil {
// 		log.Println(err)
// 		handlerError(w, CustmErr{Value: "Something Went Wrong"}, 400)
// 		return
// 	}
// 	var id = req.PostFormValue("id")
// 	println(id)
// 	var rf models.Event
// 	go func() {
// 		FindOne(d, "event", bson.M{"id": id}, &rf, &err)
// 		flagch <- 1
// 	}()
// 	if <-flagch == 1 {
// 		if err != nil && err.Error() == "not found" {
// 			handlerError(w, CustmErr{Value: "No Such Event Exist"}, 404)
// 			return
// 		}
// 		if err != nil {
// 			log.Println(err)
// 			handlerError(w, CustmErr{Value: "Something Went Wrong"}, 500)
// 			return
// 		}
// 		values := make(map[string]string)
// 		for _, v := range rf.Fields {
// 			values[v] = req.PostFormValue(v)
// 		}
// 		go func() {
// 			Update(d, "event", bson.M{"id": id}, bson.M{"$push": bson.M{"values": values}}, &err)
// 			flagch <- 2
// 		}()
// 		if <-flagch == 2 {
// 			log.Println("Switched At 2")
// 			if err != nil {
// 				log.Println(err)
// 				handlerError(w, CustmErr{Value: "Something Went Wrong"}, 500)
// 				return
// 			}
// 			bs, _ := json.Marshal(map[string]string{"msg": "Record Posted Successfully"})
// 			w.WriteHeader(200)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Write(bs)
// 		}
// 	}
// }

// func handlerError(w http.ResponseWriter, err error, sc int) {
// 	bs, _ := json.Marshal(map[string]string{"err": err.Error()})
// 	w.WriteHeader(sc)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(bs)
// }
