package router

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/flowqio/flowqlet/handlers"
)

func init() {
	log.Debug("Init Router succesful...")
	http.Handle("/", InitRouter())

}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func routerPathDebug(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log.Debugf("Path: %s, Func: %s ", path, getFuncName(f))
	r.HandleFunc(path, f)
}

//InitRouter provide single router access
func InitRouter() *mux.Router {

	r := mux.NewRouter()

	//WebSocket
	r.HandleFunc("/ws/{oid}/{cid}", handlers.ServeWS)

	// Normal page , index
	r.HandleFunc("/", handlers.Index).Methods("GET")

	//Middleware config
	//r.Use(handlers.SignCheck)
	//r.Use(handlers.GZip)

	//API Subrouter
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("", handlers.APIVersion).Methods("GET")

	//instances REST
	api.HandleFunc("/instances", handlers.GetInstance).Methods("GET")
	api.HandleFunc("/instance/{oid}", handlers.GetInstance).Methods("GET")

	api.HandleFunc("/instance/{oid}/{sid}", handlers.ComposeUP).Methods("POST")
	api.HandleFunc("/instance/{oid}/{sid}", handlers.ComposeDown).Methods("DELETE")

	//old code
	//api.HandleFunc("/instance/{oid}", handlers.CreateInstance).Methods("POST")
	//api.HandleFunc("/instance/{oid}/{cid}", handlers.ClearInstance).Methods("DELETE")

	//not found handler
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)
	return r

}
