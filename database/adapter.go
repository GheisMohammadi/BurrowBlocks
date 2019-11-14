package database

import (
	hsBC "github.com/BurrowBlocks/blockchain"
)

//UserAccount defines all user accounts
type UserAccount struct {
	ID      uint64
	Address string
	NumTxs  uint64
	Txs     []hsBC.Transaction
}

//Adapter for data base
type Adapter interface {
	Connect() error
	Disconnect() error

	//Account Handling
	InsertAccount(acc *hsBC.Account) error
	UpdateAccount(id int, acc *hsBC.Account) error
	GetAccount(id int) (*hsBC.Account, error)
	GetAccountByAddress(address string) (*hsBC.Account, error)
	GetAccountAllTransactions(address string) ([]hsBC.Transaction, error)
	GetAccountTransactions(address string,minID uint64,maxID uint64) ([]hsBC.Transaction, error)
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

	//InsertUserAccount add a unique user account in database if it not exist
	InsertUserAccount(address string, numtxs uint64) error
	//GetUserAccount returns a user account details
	GetUserAccount(address string) (*UserAccount, error)
	//UpdateUserAccount modifies all fields for selected user account
	UpdateUserAccount(address string, numtxs uint64) error
	//InsertOrAddTxToUserAccount inserts new account if not exist or add one to num_txs
	InsertOrAddTxToUserAccount(address string) error
}
