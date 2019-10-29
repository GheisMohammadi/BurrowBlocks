package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	bc "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
	mux "github.com/gorilla/mux"
)

var dbAdapter *db.Postgre
var bcAdapter *bc.Burrow
var configuration *config.Config

//Response returns response object for rest API
type Response struct {
	ErrorNumber      int                    `json:"error"`
	ErrorDescription string                 `json:"desc"`
	Result           map[string]interface{} `json:"result"`
}

//InitServer for init restful API Server
func InitServer(configObject *config.Config, dbObject *db.Postgre, bcObject *bc.Burrow) {

	configuration = configObject
	dbAdapter = dbObject
	bcAdapter = bcObject

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/info", showVersion)
	router.HandleFunc("/api/v1/gethash/{hash}", getHash).Methods("GET")
	router.HandleFunc("/api/v1/nodes", getNodesStatus).Methods("GET")

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
	res.Result = make(map[string]interface{})

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

func getNodesStatus(w http.ResponseWriter, r *http.Request) {

	var res Response
	res.Result = make(map[string]interface{})

	nodes, errGetNodesStatus := bcAdapter.GetNodes()

	if errGetNodesStatus != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get nodes: " + errGetNodesStatus.Error()
		res.Result["n_nodes"] = "0"
		res.Result["nodes"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	for _, node := range nodes {
		node.NodeInfo["listen_addr"] = "*"
		node.NodeInfo["moniker"] = "*"
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["n_nodes"] = strconv.FormatInt(int64(len(nodes)), 10)
	res.Result["nodes"] = nodes

	json.NewEncoder(w).Encode(res)
}
