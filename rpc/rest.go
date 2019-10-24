package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
	mux "github.com/gorilla/mux"
)

var dbAdapter *db.Postgre
var configuration *config.Config

//Response returns response object for rest API
type Response struct {
	ErrorNumber      int               `json:"error"`
	ErrorDescription string            `json:"desc"`
	Result           map[string]string `json:"result"`
}

//InitServer for init restful API Server
func InitServer(configObject *config.Config, dbObject *db.Postgre) {

	configuration = configObject
	dbAdapter = dbObject

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/info", showVersion)
	router.HandleFunc("/api/v1/gethash/{hash}", getHash).Methods("GET")

	//router.HandleFunc("/events", getAllEvents).Methods("GET")
	//router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	//router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	//router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	url := configuration.RestfulServer.Host + ":" + configuration.RestfulServer.Port
	http.ListenAndServe(url, router)
}

func showVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BurrowBlocks API Version 1.0 Beta")
}

func getHash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	var res Response
	res.Result = make(map[string]string)

	tx, errGetTx := dbAdapter.GetTx(hash)
	if errGetTx != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "not found"
		res.Result["height"] = "0"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["height"] = strconv.FormatInt(tx.BlockID, 10)

	json.NewEncoder(w).Encode(res)
}
