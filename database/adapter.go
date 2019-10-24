package database

import (
	hsBC "github.com/BurrowBlocks/blockchain"
)

//Adapter for data base
type Adapter interface {
	Connect() error
	Disconnect() error

	//Account Handling
	InsertAccount(acc *hsBC.Account) error
	UpdateAccount(id int, acc *hsBC.Account) error
	GetAccount(id int) (*hsBC.Account, error)
	GetAccountsTableLastID() (uint64, error)

	//Blocks Handling
	InsertBlock(b *hsBC.BlockInfo) error
	UpdateBlock(id int, b *hsBC.Block) error
	GetBlock(id int) (*hsBC.Block, error)
	GetBlocksTableLastID() (uint64, error)

	//Transactions Handling
	InsertTx(b *hsBC.Transaction) error
	UpdateTx(id int, b *hsBC.Transaction) error
	GetTx(hash string) (*hsBC.Transaction, error)
	GetTXsTableLastID() (uint64, error)

}
