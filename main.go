package main

import (
	"time"

	bc "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
	ex "github.com/BurrowBlocks/explorer"
)

var explorerEngine ex.Explorer
var gConfig *config.Config

func main() {

	Init()
	SyncLoop()

}

//Init initializes engine
func Init() {
	gConfig, _ = config.LoadConfigFile(true)
	bcAdapter := bc.Burrow{Config: gConfig}
	dbAdapter := db.Postgre{Config: gConfig}
	explorerEngine = ex.Explorer{BCAdapter: &bcAdapter, DBAdapter: &dbAdapter, Config: gConfig}

	explorerEngine.Init()
}

//SyncLoop goes in loop for syncing blockchain and database
func SyncLoop() {

	interval := time.Duration(gConfig.App.CheckingInterval)
	println("syncing every", interval, "miliseconds...")

	for {

		go func() {
			errUpdate := explorerEngine.Update()
			if errUpdate != nil {
				println("Updating engine error: ", errUpdate.Error())
			}
		}()
		time.Sleep(interval * time.Millisecond)

	}
}
