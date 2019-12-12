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
	cors "github.com/rs/cors"
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
	router.HandleFunc("/api/v1/getblock/{id}", getBlock).Methods("GET")
	router.HandleFunc("/api/v1/accountscount", getAccountsCount).Methods("GET")
	router.HandleFunc("/api/v1/accounts/{from}/{to}", getAccounts).Methods("GET")
	router.HandleFunc("/api/v1/getaccount/{address}", getAccount).Methods("GET")
	router.HandleFunc("/api/v1/getaccountalltxs/{address}", getAccountAllTxs).Methods("GET")
	router.HandleFunc("/api/v1/getaccounttxs/{address}/{minid}/{maxid}", getAccountTxs).Methods("GET")
	router.HandleFunc("/api/v1/getcumulativetxs/{barscount}", getCumulativeTxsCount).Methods("GET")
	router.HandleFunc("/api/v1/nodes", getNodesStatus).Methods("GET")
	router.HandleFunc("/api/v1/blockscount", getBlocksCount).Methods("GET")
	router.HandleFunc("/api/v1/getdurations/{count}", getDurations).Methods("GET")
	router.HandleFunc("/api/v1/txscount", getTxsCount).Methods("GET")
	router.HandleFunc("/api/v1/latesttxs/{count}", getLatestTxs).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Content-Type", "text/plain", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "ApiKey", "AdminSecretKey", "Recaptcha", "SessionToken", "Start-Row-Number", "Order-Attr", "Order-Type", "Page-Size", "Sid"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	url := configObject.RestfulServer.Host + ":" + configObject.RestfulServer.Port
	println("Rest server is listening from " + url + "...")
	http.ListenAndServe(url, handler)

}

func showVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BurrowBlocks API Version 1.0 Beta")
}

func getHash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	var res Response
	res.Result = make(map[string]interface{})

	tx, txtime, errGetTx := dbAdapter.GetTx(hash)
	if errGetTx != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "not found"
		res.Result["height"] = "0"
		res.Result["time"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["height"] = strconv.FormatInt(tx.BlockID, 10)
	res.Result["time"] = txtime

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

func getBlocksCount(w http.ResponseWriter, r *http.Request) {

	var res Response
	res.Result = make(map[string]interface{})

	count, errGetBlocksCount := dbAdapter.GetBlocksTableLastID()

	if errGetBlocksCount != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get number of saved blocks: " + errGetBlocksCount.Error()
		res.Result["num_blocks"] = "0"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["num_blocks"] = strconv.FormatUint(count, 10)

	json.NewEncoder(w).Encode(res)
}

func getBlock(w http.ResponseWriter, r *http.Request) {

	strid := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strid)

	var res Response
	res.Result = make(map[string]interface{})

	block, errGetBlock := dbAdapter.GetBlock(id)

	if errGetBlock != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get block: " + errGetBlock.Error()
		res.Result["details"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["details"] = block

	json.NewEncoder(w).Encode(res)
}

func getAccountsCount(w http.ResponseWriter, r *http.Request) {

	var res Response
	res.Result = make(map[string]interface{})

	count, errGetAccountsCount := dbAdapter.GetAccountsCount()

	if errGetAccountsCount != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get number of accounts: " + errGetAccountsCount.Error()
		res.Result["num_accs"] = "0"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["num_accs"] = strconv.FormatUint(count, 10)

	json.NewEncoder(w).Encode(res)
}

func getAccounts(w http.ResponseWriter, r *http.Request) {

	strFromID := mux.Vars(r)["from"]
	strToID := mux.Vars(r)["to"]

	fromID, _ := strconv.Atoi(strFromID)
	toID, _ := strconv.Atoi(strToID)

	var res Response
	res.Result = make(map[string]interface{})

	accs, errGetAccounts := dbAdapter.GetAccounts(uint64(fromID), uint64(toID))

	if errGetAccounts != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get accounts: " + errGetAccounts.Error()
		res.Result["num_accs"] = "0"
		res.Result["accs"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["num_accs"] = strconv.FormatUint(uint64(len(accs)), 10)
	res.Result["accs"] = accs

	json.NewEncoder(w).Encode(res)
}

func getAccount(w http.ResponseWriter, r *http.Request) {

	address := mux.Vars(r)["address"]

	var res Response
	res.Result = make(map[string]interface{})

	acc, errGetAccount := dbAdapter.GetAccountByAddress(address)

	if errGetAccount != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get account: " + errGetAccount.Error()
		res.Result["details"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["details"] = acc

	json.NewEncoder(w).Encode(res)
}

func getAccountAllTxs(w http.ResponseWriter, r *http.Request) {

	address := mux.Vars(r)["address"]

	var res Response
	res.Result = make(map[string]interface{})

	txs, errGetAccountTxs := dbAdapter.GetAccountAllTransactions(address)

	if errGetAccountTxs != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get account: " + errGetAccountTxs.Error()
		res.Result["txs"] = ""
		res.Result["totalcount"] = 0
		json.NewEncoder(w).Encode(res)
		return
	}

	if len(txs) > 0 {
		res.ErrorNumber = 0
		res.ErrorDescription = "ok"
		res.Result["txs"] = txs
		res.Result["totalcount"] = len(txs)
	} else {
		res.ErrorNumber = 1
		res.ErrorDescription = "Not Found!"
		res.Result["txs"] = ""
		res.Result["totalcount"] = 0
	}

	json.NewEncoder(w).Encode(res)
}

func getAccountTxs(w http.ResponseWriter, r *http.Request) {

	address := mux.Vars(r)["address"]
	minIDStr := mux.Vars(r)["minid"]
	maxIDStr := mux.Vars(r)["maxid"]

	s, _ := strconv.Atoi(minIDStr)
	minID := uint64(s)

	e, _ := strconv.Atoi(maxIDStr)
	maxID := uint64(e)

	var res Response
	res.Result = make(map[string]interface{})

	txs, totalCount, errGetAccountTxs := dbAdapter.GetAccountTransactions(address, minID, maxID)

	if errGetAccountTxs != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get account: " + errGetAccountTxs.Error()
		res.Result["txs"] = ""
		res.Result["totalcount"] = 0
		json.NewEncoder(w).Encode(res)
		return
	}

	if len(txs) > 0 {
		res.ErrorNumber = 0
		res.ErrorDescription = "ok"
		res.Result["txs"] = txs
		res.Result["totalcount"] = totalCount
	} else {
		res.ErrorNumber = 1
		res.ErrorDescription = "Not Found!"
		res.Result["txs"] = ""
		res.Result["totalcount"] = 0
	}

	json.NewEncoder(w).Encode(res)
}

func getCumulativeTxsCount(w http.ResponseWriter, r *http.Request) {

	strBarsCount := mux.Vars(r)["barscount"]
	ibarsCount, _ := strconv.Atoi(strBarsCount)
	barsCount := uint64(ibarsCount)

	var res Response
	res.Result = make(map[string]interface{})

	txsCountArray, errGetCumTxsCount := dbAdapter.GetCumulativeTxsCount(barsCount)

	if errGetCumTxsCount != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get cumulative number of saved txs: " + errGetCumTxsCount.Error()
		res.Result["BlockCumulativeTxs"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["BlockCumulativeTxs"] = txsCountArray

	json.NewEncoder(w).Encode(res)

}

func getDurations(w http.ResponseWriter, r *http.Request) {

	strCount := mux.Vars(r)["count"]
	iCount, _ := strconv.Atoi(strCount)
	Count := uint64(iCount)

	var res Response
	res.Result = make(map[string]interface{})

	durations, errGetDurations := dbAdapter.GetBlocksDurations(Count)

	if errGetDurations != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get durations of block: " + errGetDurations.Error()
		res.Result["durations"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["durations"] = durations

	json.NewEncoder(w).Encode(res)

}

func getTxsCount(w http.ResponseWriter, r *http.Request) {

	var res Response
	res.Result = make(map[string]interface{})

	count, errGetTxsCount := dbAdapter.GetTxsCount()

	if errGetTxsCount != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get number of saved txs: " + errGetTxsCount.Error()
		res.Result["num_txs"] = "0"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.ErrorNumber = 0
	res.ErrorDescription = "ok"
	res.Result["num_txs"] = strconv.FormatUint(count, 10)

	json.NewEncoder(w).Encode(res)
}

func getLatestTxs(w http.ResponseWriter, r *http.Request) {

	strcount := mux.Vars(r)["count"]
	c, _ := strconv.Atoi(strcount)
	count := uint64(c)

	var res Response
	res.Result = make(map[string]interface{})

	txs, errGetLatestTxs := dbAdapter.GetLatestTxs(count)

	if errGetLatestTxs != nil {
		res.ErrorNumber = 1
		res.ErrorDescription = "can't get latest txs: " + errGetLatestTxs.Error()
		res.Result["txs"] = ""
		json.NewEncoder(w).Encode(res)
		return
	}

	if len(txs) > 0 {
		res.ErrorNumber = 0
		res.ErrorDescription = "ok"
		res.Result["txs"] = txs
	} else {
		res.ErrorNumber = 1
		res.ErrorDescription = "Not Found!"
		res.Result["txs"] = ""
	}

	json.NewEncoder(w).Encode(res)
}
