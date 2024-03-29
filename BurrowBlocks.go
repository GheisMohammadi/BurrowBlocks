package main

import (
	"time"

	bc "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
	ex "github.com/BurrowBlocks/explorer"

	rest "github.com/BurrowBlocks/rpc"
)

var explorerEngine ex.Explorer
var dbAdapter db.Postgre
var bcAdapter bc.Burrow
var gConfig *config.Config

func main() {

	//Initializing...
	Init()

	//Prepairing Restful API...
	go func() {
		//defer dbAdapter.Disconnect()
		rest.InitServer(gConfig, &dbAdapter, &bcAdapter)
	}()

	//Sync with BlockChain...
	SyncLoop()

}

//Init initializes engine
func Init() {
	gConfig, _ = config.LoadConfigFile(true)
	bcAdapter = bc.Burrow{Config: gConfig}
	dbAdapter = db.Postgre{Config: gConfig}
	explorerEngine = ex.Explorer{BCAdapter: &bcAdapter, DBAdapter: &dbAdapter, Config: gConfig}

	explorerEngine.Init()
}

//SyncLoop goes in loop for syncing blockchain and database
func SyncLoop() {

	interval := time.Duration(gConfig.App.CheckingInterval)
	println("syncing every", interval, "miliseconds...")
	defer dbAdapter.Disconnect()

	for {

		//go func() {
		errUpdate := explorerEngine.Update()
		if errUpdate != nil {
			println("Updating engine error: ", errUpdate.Error())
		}
		//}()
		time.Sleep(interval * time.Millisecond)

	}
}
