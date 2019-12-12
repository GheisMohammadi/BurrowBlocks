package explorer

import (
	"fmt"

	bc "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
)

var numDots int

//Explorer class for connecting block chain to data base
type Explorer struct {
	BCAdapter bc.Adapter
	DBAdapter db.Adapter
	Config    *config.Config
}

//Init to initialize database and block chain
func (e *Explorer) Init() error {
	//connect to gallactic blockchain by gRPC
	bcAdapter := e.BCAdapter
	clientErr := bcAdapter.CreateClient()
	if clientErr == nil {
		println("Blockchain client created successfully!")
	} else {
		return clientErr
	}
	bcAdapter.Update()

	//connect to database
	dbAdapter := e.DBAdapter
	connErr := dbAdapter.Connect()
	if connErr != nil {
		return connErr
	}
	println("Connected to database successfully!")

	return nil
}

//Update to sync
func (e *Explorer) Update() error {
	//update db if need
	return e.UpdateAll()
}

//UpdateAll to Sync database with blockchain
func (e *Explorer) UpdateAll() error {

	//Sync data with blockchain
	updateErr := e.BCAdapter.Update()
	if updateErr != nil {
		return updateErr
	}

	//Get current height of last block
	currentHeight, getLastHeightErr := e.BCAdapter.GetBlocksLastHeight()
	if getLastHeightErr != nil {
		println("error on reading last block id: " + getLastHeightErr.Error())
		return getLastHeightErr
	}

	//Get last block ID that is saved
	lastBlockIDInDB, getLastIDError := e.DBAdapter.GetBlocksTableLastID()
	if getLastIDError != nil {
		println("error on reading last block id from db: " + getLastIDError.Error())
		lastBlockIDInDB = 0
	}

	/*
		inf,errGetBlockInfo := bcAdapter.GetBlockInfo(8)
		if errGetBlockInfo!=nil{
			println("get block info error: ",errGetBlockInfo)
		}
		println("Block Info: ", inf.TotalTxs)
	*/

	if currentHeight > lastBlockIDInDB {
		d := currentHeight - lastBlockIDInDB
		n := int(d / 1000)
		if d > 1000 {
			s := lastBlockIDInDB
			if s == 0 {
				s = 1
			}
			println("Saving new blocks number", s, " to ", currentHeight, " in database...")
		}

		var startIndex uint64
		var endIndex uint64

		startBlockID := lastBlockIDInDB + 1

		for i := 0; i <= n; i++ {
			startIndex = startBlockID + uint64(i*1000)
			endIndex = startIndex + 999
			if endIndex > currentHeight {
				endIndex = currentHeight
			}

			blocks, _ := e.BCAdapter.GetBlocks(startIndex, endIndex)
			savingErr := e.saveBlocksInDB(blocks, e.BCAdapter, e.DBAdapter)
			if savingErr != nil {
				println("error on saving blocks in db: " + savingErr.Error())
				return savingErr
			}

			perc := (int)((float64(i+1) / float64(n+1)) * 100.0)
			fmt.Printf("\r%d%% saved! (%d/%d)", perc, endIndex-startBlockID, d)
		}

		if d > 1000 {
			println("\r", d, "new blocks saved!                   ")
			println("Checking new blocks...")
		} else {
			e.writeAnim(currentHeight)
		}

	} else {
		e.writeAnim(currentHeight)
	}

	return nil
}

func (e *Explorer) saveBlocksInDB(blocks []bc.BlockInfo, bcAdapter bc.Adapter, dbAdapter db.Adapter) error {
	l := len(blocks)
	if l <= 0 {
		return fmt.Errorf("Empty Blocks Array")
	}
	for i := 0; i < l; i++ {
		block := blocks[i]
		err := dbAdapter.InsertBlock(&block)
		if err != nil {
			println("error on insert block in db: " + err.Error())
			return err
		}
		if block.NumTxs > 0 {
			errTxSave := e.saveBlockTXsInDB(block, bcAdapter, dbAdapter)
			if errTxSave != nil {
				println("error on save block txs in db: " + errTxSave.Error())
				return errTxSave
			}
		}
	}

	syncInfo, errGetSyncInfo := bcAdapter.GetSyncInfo()
	if errGetSyncInfo != nil {
		println("error on get sync info: " + errGetSyncInfo.Error())
		return errGetSyncInfo
	}

	blockID := int64(syncInfo.LatestBlockHeight)
	if blockID==blocks[l-1].Height {
		duration := syncInfo.LatestBlockDuration
		errUpdateDuration := dbAdapter.UpdateBlockDuration(blockID, duration)
		if errUpdateDuration != nil {
			println("error on update block duration: " + errUpdateDuration.Error())
			return errUpdateDuration
		}
	}

	return nil
}

func (e *Explorer) saveBlockTXsInDB(block bc.BlockInfo, bcAdapter bc.Adapter, dbAdapter db.Adapter) error {
	l := block.NumTxs
	if l <= 0 {
		return fmt.Errorf("Empty Transactions Array")
	}

	height := uint64(block.Height)
	txs, errTXs := bcAdapter.GetTXs(height)
	if errTXs != nil {
		println("error on get block txs: " + errTXs.Error())
		return errTXs
	}

	if int64(len(txs)) != l {
		return fmt.Errorf("error on parsing txs for block %v some txs are missed", height)
	}

	for i := l - 1; i >= 0; i-- {
		err := dbAdapter.InsertTx(&txs[i])
		if err != nil {
			println("error on save tx in db: " + err.Error())
			return err
		}
	}
	return nil
}

func (e *Explorer) writeAnim(currentHeight uint64) {
	numDots++
	if numDots > 3 {
		numDots = 0
	}
	dotStr := ""
	for nDot := 1; nDot <= numDots; nDot++ {
		dotStr += "."
	}
	fmt.Printf("\r%d blocks saved"+dotStr+"       ", currentHeight)
}
