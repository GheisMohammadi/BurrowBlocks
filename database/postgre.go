package database

import (
	"database/sql"
	"fmt"

	hsBC "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	_ "github.com/lib/pq" //dependency for postgre
)

//Postgre adapter
type Postgre struct {
	Config *config.Config
	ObjDB  *sql.DB //Opened DB
}

//Connect to database
func (obe *Postgre) Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		obe.Config.DataBase.Host, obe.Config.DataBase.Port, obe.Config.DataBase.User,
		obe.Config.DataBase.Password, obe.Config.DataBase.DBName)

	var err error
	obe.ObjDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		obe.ObjDB.Close()
		return err
	}

	err = obe.ObjDB.Ping()
	if err != nil {
		obe.ObjDB.Close()
		return err
	}
	return nil
}

//Disconnect close connection to database
func (obe *Postgre) Disconnect() error {
	closeError := obe.ObjDB.Close()
	return closeError
}

//InsertAccount add new Account to accounts table
func (obe *Postgre) InsertAccount(acc *hsBC.Account) error {

	sqlStatement := `INSERT INTO accounts (address, balance, permission,sequence,code)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`
	id := 0
	err := obe.ObjDB.QueryRow(sqlStatement, acc.Address, acc.Balance, acc.Permission, acc.Sequence, acc.Code).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateAccount modifies all fields for selected account
func (obe *Postgre) UpdateAccount(id int, acc *hsBC.Account) error {
	sqlStatement := `UPDATE accounts
				SET address = $2, balance = $3, permission = $4, sequence = $5, code = $6
				WHERE id = $1
				RETURNING id, address;`
	var retAddress string
	var retID int
	err := obe.ObjDB.QueryRow(sqlStatement, id, acc.Address, acc.Balance, acc.Permission, acc.Sequence, acc.Code).Scan(&retID, &retAddress)

	if err != nil {
		return err
	}

	return nil
}

//GetAccount finds account in db and returns its data
func (obe *Postgre) GetAccount(id int) (*hsBC.Account, error) {
	sqlStatement := `SELECT * FROM accounts 
					 WHERE id=$1;`
	acc := &hsBC.Account{Address: "", Balance: 0.0, Permission: "", Sequence: 0, Code: ""}
	row := obe.ObjDB.QueryRow(sqlStatement, id)
	err := row.Scan(&acc.Address, &acc.ID, &acc.Balance, &acc.Permission, &acc.Sequence, &acc.Code)
	switch err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		return acc, nil
	default:
		return nil, err
	}
}

//GetAccountByAddress finds account in db and returns its data
func (obe *Postgre) GetAccountByAddress(address string) (*hsBC.Account, error) {
	sqlStatement := `SELECT * FROM accounts 
					 WHERE address=$1;`
	acc := &hsBC.Account{Address: "", Balance: 0.0, Permission: "", Sequence: 0, Code: ""}
	row := obe.ObjDB.QueryRow(sqlStatement, address)
	err := row.Scan(&acc.Address, &acc.ID, &acc.Balance, &acc.Permission, &acc.Sequence, &acc.Code)
	switch err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		return acc, nil
	default:
		return nil, err
	}
}

//GetAccountFromTransactions finds accounts in db.transactions using address and returns its data
func (obe *Postgre) GetAccountFromTransactions(address string) ([]hsBC.Transaction, error) {
	sqlStatement := `SELECT block_id,txhash,fee,gas_limit,data,addr_from,addr_to,amount,tx_type FROM transactions 
					 WHERE addr_from=$1 OR addr_to=$1;`

	rows, errGetTxs := obe.ObjDB.Query(sqlStatement, address)
	defer rows.Close()

	if errGetTxs != nil {
		return nil, errGetTxs
	}

	txs := make([]hsBC.Transaction, 0)
	for rows.Next() {

		var txn hsBC.Transaction
		if err := rows.Scan(&txn.BlockID, &txn.Hash, &txn.Fee, &txn.GasLimit, &txn.Data, &txn.From, &txn.To, &txn.Amount, &txn.Type); err != nil {
			return nil, err
		}

		txs = append(txs, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return txs, nil
}

//GetAccountsTableLastID returns last block number
func (obe *Postgre) GetAccountsTableLastID() (uint64, error) {
	sqlStatement := `SELECT coalesce(MAX(id), 0) as max FROM accounts;`

	row := obe.ObjDB.QueryRow(sqlStatement)
	var LastID uint64
	err := row.Scan(&LastID)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return LastID, nil
	default:
		return 0, err
	}
}

//InsertBlock add a block in database
func (obe *Postgre) InsertBlock(b *hsBC.BlockInfo) error {
	sqlStatement := `INSERT INTO blocks (height, hash, chainID, time, txcounts)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING height`
	id := 0
	row := obe.ObjDB.QueryRow(sqlStatement, b.Height, b.BlockHash, b.ChainID, b.Time, b.NumTxs)
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateBlock modifies a block data in database
func (obe *Postgre) UpdateBlock(id int, b *hsBC.Block) error {
	//TODO:
	return nil
}

//GetBlock returns a block details
func (obe *Postgre) GetBlock(id int) (*hsBC.Block, error) {
	sqlStatement := `SELECT height, hash, chainID, time, txcounts FROM blocks 
					 WHERE height=$1;`

	var b hsBC.Block
	row := obe.ObjDB.QueryRow(sqlStatement, id)
	err := row.Scan(&b.Height, &b.Hash, &b.ChainID, &b.Time, &b.TxCounts)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

//GetBlocksTableLastID returns last block number
func (obe *Postgre) GetBlocksTableLastID() (uint64, error) {
	sqlStatement := `SELECT coalesce(MAX(height), 0) as max FROM blocks`

	row := obe.ObjDB.QueryRow(sqlStatement)

	var LastID uint64
	err := row.Scan(&LastID)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return LastID, nil
	default:
		return 0, err
	}
}

//GetBlocksCount returns num blocks saved in db
func (obe *Postgre) GetBlocksCount() (uint64, error) {
	sqlStatement := `SELECT COUNT(*) as count FROM blocks`

	row := obe.ObjDB.QueryRow(sqlStatement)
	var TxsCount uint64
	err := row.Scan(&TxsCount)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return TxsCount, nil
	default:
		return 0, err
	}
}

//InsertTx add a transaction in database
func (obe *Postgre) InsertTx(b *hsBC.Transaction) error {
	sqlStatement := `INSERT INTO transactions (block_id, txhash, fee, gas_limit, data, addr_from, addr_to, amount, tx_type)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`
	id := 0
	row := obe.ObjDB.QueryRow(sqlStatement, b.BlockID, b.Hash, b.Fee, b.GasLimit, b.Data, b.From, b.To, b.Amount, b.Type)
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	if b.From != "" {
		err = obe.InsertOrAddTxToUserAccount(b.From)
		if err != nil {
			return err
		}
	}

	if b.To != "" {
		err = obe.InsertOrAddTxToUserAccount(b.To)
		if err != nil {
			return err
		}
	}

	return nil
}

//UpdateTx modifies a transaction data in database
func (obe *Postgre) UpdateTx(id int, b *hsBC.Transaction) error {
	sqlStatement := `UPDATE transactions
	SET block_id = $2, txhash = $3, fee = $4, gas_limit = $5, data = $6, addr_from = $7, addr_to = $8, amount = $9, tx_type = $10
	WHERE id = $1
	RETURNING id, txhash;`
	var retHash string
	var retID int
	err := obe.ObjDB.QueryRow(sqlStatement, id, b.BlockID, b.Hash, b.Fee, b.GasLimit, b.Data, b.From, b.To, b.Amount, b.Type).Scan(&retID, &retHash)

	if err != nil {
		return err
	}

	return nil
}

//GetTx returns a transaction data
func (obe *Postgre) GetTx(hash string) (*hsBC.Transaction, error) {
	sqlStatement := `SELECT block_id,txhash,fee,gas_limit,data,addr_from,addr_to,amount,tx_type FROM transactions
					 WHERE txhash=$1;`
	var tx hsBC.Transaction
	errGetTx := obe.ObjDB.QueryRow(sqlStatement, hash).Scan(&tx.BlockID, &tx.Hash, &tx.Fee, &tx.GasLimit, &tx.Data, &tx.From, &tx.To, &tx.Amount, &tx.Type)

	if errGetTx != nil {
		return nil, errGetTx
	}

	return &tx, nil
}

//GetTXsTableLastID returns last saved transaction number
func (obe *Postgre) GetTXsTableLastID() (uint64, error) {
	sqlStatement := `SELECT coalesce(MAX(id), 0) as max FROM transactions`

	row := obe.ObjDB.QueryRow(sqlStatement)
	var LastID uint64
	err := row.Scan(&LastID)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return LastID, nil
	default:
		return 0, err
	}
}

//GetTxsCount returns num transaction saved in db
func (obe *Postgre) GetTxsCount() (uint64, error) {
	sqlStatement := `SELECT COUNT(*) as count FROM transactions`

	row := obe.ObjDB.QueryRow(sqlStatement)
	var TxsCount uint64
	err := row.Scan(&TxsCount)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return TxsCount, nil
	default:
		return 0, err
	}
}

//GetLatestTxs returns latest transactions by count
func (obe *Postgre) GetLatestTxs(count uint64) ([]hsBC.Transaction, error) {

	sqlStatement := `SELECT block_id,txhash,fee,gas_limit,data,addr_from,addr_to,amount,tx_type FROM transactions
	ORDER BY block_id DESC LIMIT $1;`

	rows, errGetLatestTxs := obe.ObjDB.Query(sqlStatement, count)
	defer rows.Close()

	if errGetLatestTxs != nil {
		return nil, errGetLatestTxs
	}

	txs := make([]hsBC.Transaction, 0)
	for rows.Next() {

		var txn hsBC.Transaction
		if err := rows.Scan(&txn.BlockID, &txn.Hash, &txn.Fee, &txn.GasLimit, &txn.Data, &txn.From, &txn.To, &txn.Amount, &txn.Type); err != nil {
			return nil, err
		}

		txs = append(txs, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return txs, nil
}

//GetAccountsCount returns number of user accounts
func (obe *Postgre) GetAccountsCount() (uint64, error) {

	sqlStatement := `SELECT COUNT(*) as count FROM useraccounts`

	row := obe.ObjDB.QueryRow(sqlStatement)
	var UserAccsCount uint64
	err := row.Scan(&UserAccsCount)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return UserAccsCount, nil
	default:
		return 0, err
	}
}

//GetAccounts returns all account of users
func (obe *Postgre) GetAccounts(fromID uint64, toID uint64) ([]UserAccount, error) {

	sqlStatement := `SELECT * FROM useraccounts 
	WHERE id>=$1 AND id<=$2;`

	rows, errGetUserAccs := obe.ObjDB.Query(sqlStatement, fromID, toID)
	//err := row.Scan(&acc.Address, &acc.ID, &acc.Address, &acc.NumTxs)

	defer rows.Close()

	if errGetUserAccs != nil {
		return nil, errGetUserAccs
	}

	accs := make([]UserAccount, 0)
	for rows.Next() {

		var acc UserAccount
		if err := rows.Scan(&acc.ID, &acc.Address, &acc.NumTxs); err != nil {
			return nil, err
		}

		accs = append(accs, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accs, nil
}

//InsertUserAccount add a unique user account in database if it not exist
func (obe *Postgre) InsertUserAccount(address string, numtxs uint64) error {

	sqlStatement := `INSERT INTO useraccounts (address, num_txs)
	VALUES ($1, $2)
	RETURNING id`

	id := 0
	row := obe.ObjDB.QueryRow(sqlStatement, address, numtxs)
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

//GetUserAccount returns a user account details
func (obe *Postgre) GetUserAccount(address string) (*UserAccount, error) {
	sqlStatement := `SELECT id, address, num_txs FROM useraccounts
					 WHERE address=$1;`
	var acc UserAccount
	errGetAcc := obe.ObjDB.QueryRow(sqlStatement, address).Scan(&acc.ID, &acc.Address, &acc.NumTxs)

	if errGetAcc != nil {
		return nil, errGetAcc
	}

	return &acc, nil
}

//UpdateUserAccount modifies all fields for selected user account
func (obe *Postgre) UpdateUserAccount(address string, numtxs uint64) error {
	sqlStatement := `UPDATE useraccounts
				SET num_txs = $2
				WHERE address = $1
				RETURNING id;`
	var retID int
	err := obe.ObjDB.QueryRow(sqlStatement, address, numtxs).Scan(&retID)

	if err != nil {
		return err
	}

	return nil
}

//InsertOrAddTxToUserAccount inserts new account if not exist or add one to num_txs
func (obe *Postgre) InsertOrAddTxToUserAccount(address string) error {

	acc, errGetAcc := obe.GetUserAccount(address)
	if errGetAcc != nil && errGetAcc != sql.ErrNoRows {
		return errGetAcc
	}

	var err error
	if errGetAcc == sql.ErrNoRows {
		err = obe.InsertUserAccount(address, 1)
	} else {
		err = obe.UpdateUserAccount(address, acc.NumTxs+1)
	}

	return err
}
