package blockchain

import (
	"time"
)

//Account struct
type Account struct {
	Address    string
	Balance    uint64
	Permission string
	Sequence   uint64
	Code       string
	ID         uint64
}

//BlockInfo struct in blocks
type BlockInfo struct {
	//block ID
	BlockHash string
	// basic block info
	//VersionBlock uint64
	//VersionApp   uint64
	ChainID  string
	Height   int64
	Time     string
	NumTxs   int64
	TotalTxs int64
	// prev block info
	//LastBlockHash string
	// hashes of block data
	LastCommitHash string
	DataHash       string
	// hashes from the app output from the prev block
	ValidatorsHash     string
	NextValidatorsHash string
	ConsensusHash      string
	AppHash            string
	LastResultsHash    string
	// consensus info
	EvidenceHash    string
	ProposerAddress string
}

//Block struct
type Block struct {
	Height  int64
	Hash    string
	ChainID string
	Time    time.Time
	//LastBlockHash string
	TxCounts int64
}

//Transaction struct
type Transaction struct {
	Type     string
	BlockID  int64
	Hash     string
	From     string
	To       string
	Amount   uint64
	GasLimit uint64
	Fee      uint64
	Data     string
	Sequence uint64
}

//NodeConnectionStatus for get status of connection
type NodeConnectionStatus struct {
	Duration    int64                    `json:"Duration"`
	SendMonitor map[string]interface{}   `json:"SendMonitor"`
	RecvMonitor map[string]interface{}   `json:"RecvMonitor"`
	Channels    []map[string]interface{} `json:"Channels"`
	RemoteIP    string                   `json:"remote_ip"`
}

//Peer in network response for staus of node
type Peer struct {
	NodeInfo         map[string]interface{} `json:"node_info"`
	IsOutbound       bool                   `json:"is_outbound"`
	ConnectionStatus NodeConnectionStatus   `json:"connection_status"`
	RemoteIP         string                 `json:"remote_ip"`
}

//Adapter for data base
type Adapter interface {
	CreateClient() error

	Update() error

	GetAccountsCount() int
	GetAccount(id int) (*Account, error)
	GetAccounts() ([]*Account, error)

	GetBlocksLastHeight() (uint64, error)
	GetBlockInfo(height uint64) (*BlockInfo, error)
	GetBlock(height uint64) (*Block, error)
	GetBlocksInfo(from uint64, to uint64) ([]BlockInfo, error)
	GetBlocks(from uint64, to uint64) ([]BlockInfo, error)

	GetTXsCount(height uint64) int
	GetTx(height uint64, hash []byte) (*Transaction, error)
	GetTXs(height uint64) ([]Transaction, error)

	GetNodes() ([]Peer, error)
}
