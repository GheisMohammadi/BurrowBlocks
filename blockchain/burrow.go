package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	//"strings"

	config "github.com/BurrowBlocks/config"
)

//Burrow class for connecting to Gallactic block chain
type Burrow struct {
	Config *config.Config
	Domain string
	tr     *http.Transport
	client *http.Client
}

//CreateGRPCClient creates a client for communicating with gallactic blockchain
func (g *Burrow) CreateGRPCClient() error {
	var connURL string
	connURL = g.Config.GRPC.URL + ":" + g.Config.GRPC.Port
	g.Domain = connURL

	// Also add timeouts for connections
	g.tr = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		//DisableKeepAlives:     false,
	}

	g.client = &http.Client{
		Transport: g.tr,
		Timeout:   time.Second * 5,
	}

	return nil
}

//Update will refresh all data and sync with block chain
func (g *Burrow) Update() error {
	//var chainInfoErr error

	/*
		_, chainInfoErr = http.Get(g.Domain + "/chain_id")
		if chainInfoErr != nil {
			println("error on reading chain id info: " + chainInfoErr.Error())
			return chainInfoErr
		}
	*/

	//return chainInfoErr

	return nil
}

//GetReply for issuing get to remote url
func (g *Burrow) GetReply(strURL string) string {

	defer g.tr.CloseIdleConnections()
	defer g.client.CloseIdleConnections()

	// Turn it into a request
	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		fmt.Println("\nError forming request: " + err.Error())
		return ""
	}

	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	req.Close = true

	// Get the URL
	res, err := g.client.Do(req)
	if err != nil {
		fmt.Println("\nError reading response body: " + err.Error())
		if res != nil {
			res.Body.Close()
		}
		return ""
	}

	res.Close = true

	// What was the response status from the server?
	var strResult string
	if res.StatusCode != 200 {
		fmt.Println("\nError reading response body, status code: " + res.Status)
		if res != nil {
			res.Body.Close()
		}
		return ""
	}

	// Read the reply
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("\nError reading response body: " + err.Error())
		if res != nil {
			res.Body.Close()
		}
		return ""
	}
	res.Body.Close()

	// Cut down on calls to convert this
	strResult = string(body)

	// Done
	return strResult

}

//GetBlocksLastHeight return last height
func (g *Burrow) GetBlocksLastHeight() (uint64, error) {

	/*
		response, consInfoErr := http.Get(g.Domain + "/consensus")

		if consInfoErr != nil {
			println("get height error: ", consInfoErr)
			return 0, consInfoErr
		}
		response.Close = true
		defer response.Body.Close()

		responseData, _ := ioutil.ReadAll(response.Body)
	*/

	responseData := g.GetReply(g.Domain + "/consensus")

	type Result struct {
		RoundState map[string]interface{}   `json:"round_state"`
		Peers      []map[string]interface{} `json:"peers"`
	}

	type Consensus struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  Result `json:"result"`
	}

	//var env JsonMsg
	var env Consensus
	bytes := []byte(string(responseData))
	if err := json.Unmarshal(bytes, &env); err != nil {
		return 0, err
	}

	heightStr := fmt.Sprintf("%v", env.Result.RoundState["height"])
	height, _ := strconv.ParseUint(heightStr, 10, 64)

	return height - 1, nil
}

//GetBlockInfo returns specified block
func (g *Burrow) GetBlockInfo(height uint64) (*BlockInfo, error) {

	var inf BlockInfo

	url := fmt.Sprintf(g.Domain+"/block?height=%v", height)
	/*
		response, blockInfoErr := http.Get(url)

		if blockInfoErr != nil {
			println("get block info response error: ", blockInfoErr.Error())
			return nil, blockInfoErr
		}
		response.Close = true
		defer response.Body.Close()

		responseData, _ := ioutil.ReadAll(response.Body)
	*/

	responseData := g.GetReply(url)

	type BlockID struct {
		Hash  string                 `json:"hash"`
		Parts map[string]interface{} `json:"parts"`
	}

	type BlockMetaDetails struct {
		ID     BlockID                `json:"block_id"`
		Header map[string]interface{} `json:"header"`
	}

	type BlockData struct {
		Header     map[string]interface{} `json:"header"`
		Data       map[string]interface{} `json:"data"`
		Evidence   map[string]interface{} `json:"evidence"`
		LastCommit map[string]interface{} `json:"last_commit"`
	}

	type Result struct {
		BlockMeta BlockMetaDetails `json:"BlockMeta"`
		Block     BlockData        `json:"Block"`
	}

	type BlockDetails struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  Result `json:"result"`
	}

	//var env JsonMsg
	var env BlockDetails
	bytes := []byte(string(responseData))
	if err := json.Unmarshal(bytes, &env); err != nil {
		return nil, err
	}

	heightStr := fmt.Sprintf("%v", env.Result.Block.Header["height"])
	blockHeight, _ := strconv.ParseInt(heightStr, 10, 64)

	numTxsStr := fmt.Sprintf("%v", env.Result.Block.Header["num_txs"])
	numTxs, _ := strconv.ParseInt(numTxsStr, 10, 64)

	totalTxsStr := fmt.Sprintf("%v", env.Result.Block.Header["total_txs"])
	totalTxs, _ := strconv.ParseInt(totalTxsStr, 10, 64)

	//inf.VersionBlock = 0 //fmt.Sprintf("%v", env.Result.Block.Header["version"])
	//inf.VersionApp = 0
	inf.ChainID = fmt.Sprintf("%v", env.Result.Block.Header["chain_id"])
	inf.Height = blockHeight
	inf.Time = fmt.Sprintf("%v", env.Result.Block.Header["time"])
	inf.NumTxs = numTxs
	inf.TotalTxs = totalTxs
	inf.BlockHash = fmt.Sprintf("%v", env.Result.BlockMeta.ID.Hash)
	inf.LastBlockHash = "" //fmt.Sprintf("%v", env.Result.Block.Header[""])
	inf.DataHash = fmt.Sprintf("%v", env.Result.Block.Header["data_hash"])
	inf.ValidatorsHash = fmt.Sprintf("%v", env.Result.Block.Header["validators_hash"])
	inf.NextValidatorsHash = fmt.Sprintf("%v", env.Result.Block.Header["next_validators_hash"])
	inf.ConsensusHash = fmt.Sprintf("%v", env.Result.Block.Header["consensus_hash"])
	inf.AppHash = fmt.Sprintf("%v", env.Result.Block.Header["app_hash"])
	inf.LastResultsHash = fmt.Sprintf("%v", env.Result.Block.Header["last_results_hash"])
	inf.EvidenceHash = fmt.Sprintf("%v", env.Result.Block.Header["evidence_hash"])
	inf.ProposerAddress = fmt.Sprintf("%v", env.Result.Block.Header["proposer_address"])

	return &inf, nil
}

//GetBlock returns specified block
func (g *Burrow) GetBlock(height uint64) (*Block, error) {

	return nil, nil

	/*
		lastID, lastIDErr := g.GetBlocksLastHeight()
		if lastIDErr != nil {
			return nil, lastIDErr
		}
		if height > lastID {
			return nil, fmt.Errorf("block height out of range (max is " + string(lastID) + ")")
		}

		client := *g.client
		blockRes, getBlockErr := client.GetBlock(context.Background(), &pb.BlockRequest{Height: height})
		if getBlockErr != nil {
			return nil, getBlockErr
		}
		var b Block
		toBlock(blockRes, &b)

		return &b, nil
	*/
}

//GetAccountsCount returns number of accounts
func (g *Burrow) GetAccountsCount() int {
	return 0
	/*
		l := len(g.accounts.Accounts)
		return l
	*/
}

//GetAccount returns specified account
func (g *Burrow) GetAccount(id int) (*Account, error) {
	return nil, nil
	/*
		acc := g.accounts.Accounts[id].Account
		ID := uint64(id)
		var retAcc Account
		toAccount(acc, &retAcc)
		retAcc.ID = ID
		return &retAcc, nil
	*/
}

//GetAccounts returns all accounts in array of accounts
func (g *Burrow) GetAccounts() ([]*Account, error) {
	return nil, nil
	/*
		l := len(g.accounts.Accounts)

		retAccounts := make([]*Account, l)

		for i := 0; i < l; i++ {
			acc := g.accounts.Accounts[i].Account
			ID := uint64(i)
			toAccount(acc, retAccounts[i])
			retAccounts[i].ID = ID
		}

		return retAccounts, nil
	*/
}

//GetBlocksInfo returns a group of blocks for faster access them
func (g *Burrow) GetBlocksInfo(from uint64, to uint64) ([]BlockInfo, error) {
	return nil, nil
	/*
		client := *g.client
		blocks, getBlocksErr := client.GetBlocks(context.Background(), &pb.BlocksRequest{MinHeight: from, MaxHeight: to})
		if getBlocksErr != nil {
			return nil, getBlocksErr
		}

		n := len(blocks.GetBlocks())
		retBlocks := make([]*BlockInfo, 0, n)
		for i := 0; i < n; i++ {
			toBlockInfo(blocks.GetBlocks()[i].GetHeader(), retBlocks[i])
		}

		return retBlocks, nil
	*/
}

//GetBlocks returns a group of blocks for faster access them
func (g *Burrow) GetBlocks(from uint64, to uint64) ([]BlockInfo, error) {

	url := fmt.Sprintf(g.Domain+"/blocks?minHeight=%v&maxHeight=%v", from, to)
	/*
		response, blockInfoErr := http.Get(url)

		if blockInfoErr != nil {
			println("get blocks response error: ", blockInfoErr.Error())
			return nil, blockInfoErr
		}
		response.Close = true
		defer response.Body.Close()

		responseData, _ := ioutil.ReadAll(response.Body)
	*/

	responseData := g.GetReply(url)

	type BlockID struct {
		Hash  string                 `json:"hash"`
		Parts map[string]interface{} `json:"parts"`
	}

	type BlockMetaDetails struct {
		ID     BlockID                `json:"block_id"`
		Header map[string]interface{} `json:"header"`
	}

	type Result struct {
		LastHeight uint64             `json:"LastHeight"`
		BlockMetas []BlockMetaDetails `json:"BlockMetas"`
	}

	type BlockDetails struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  Result `json:"result"`
	}

	//var env JsonMsg
	var env BlockDetails
	bytes := []byte(string(responseData))
	if err := json.Unmarshal(bytes, &env); err != nil {
		println("unmarshal block data error: ", err.Error())
		return nil, err
	}

	if len(env.Result.BlockMetas) <= 0 {
		println("empty blocks error")
		return nil, fmt.Errorf("no blocks exist in this range")
	}

	//nBlocks := len(env.Result.BlockMetas)
	blocks := make([]BlockInfo, 0)

	for _, meta := range env.Result.BlockMetas {

		heightStr := fmt.Sprintf("%v", meta.Header["height"])
		blockHeight, _ := strconv.ParseInt(heightStr, 10, 64)

		numTxsStr := fmt.Sprintf("%v", meta.Header["num_txs"])
		numTxs, _ := strconv.ParseInt(numTxsStr, 10, 64)

		totalTxsStr := fmt.Sprintf("%v", meta.Header["total_txs"])
		totalTxs, _ := strconv.ParseInt(totalTxsStr, 10, 64)

		var inf BlockInfo

		//inf.VersionBlock = 0 //fmt.Sprintf("%v", env.Result.Block.Header["version"])
		//inf.VersionApp = 0
		inf.ChainID = fmt.Sprintf("%v", meta.Header["chain_id"])
		inf.Height = blockHeight
		inf.Time = fmt.Sprintf("%v", meta.Header["time"])
		inf.NumTxs = numTxs
		inf.TotalTxs = totalTxs
		inf.BlockHash = fmt.Sprintf("%v", meta.ID.Hash)
		inf.LastBlockHash = "" //fmt.Sprintf("%v", meta.Header[""])
		inf.DataHash = fmt.Sprintf("%v", meta.Header["data_hash"])
		inf.ValidatorsHash = fmt.Sprintf("%v", meta.Header["validators_hash"])
		inf.NextValidatorsHash = fmt.Sprintf("%v", meta.Header["next_validators_hash"])
		inf.ConsensusHash = fmt.Sprintf("%v", meta.Header["consensus_hash"])
		inf.AppHash = fmt.Sprintf("%v", meta.Header["app_hash"])
		inf.LastResultsHash = fmt.Sprintf("%v", meta.Header["last_results_hash"])
		inf.EvidenceHash = fmt.Sprintf("%v", meta.Header["evidence_hash"])
		inf.ProposerAddress = fmt.Sprintf("%v", meta.Header["proposer_address"])

		//blocks[i] = &inf
		blocks = append(blocks, inf)
	}

	return blocks, nil
}

/*
type BlockTxsResponse struct {
	Count                int32                                         `protobuf:"varint,1,opt,name=Count,proto3" json:"Count,omitempty"`
	Txs                  []github_com_gallactic_gallactic_txs.Envelope `protobuf:"bytes,3,rep,name=Txs,proto3,customtype=github.com/gallactic/gallactic/txs.Envelope" json:"Txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

type Envelope struct {
	ChainID     string             `json:"chainId"`
	Type        tx.Type            `json:"type"`
	Tx          tx.Tx              `json:"tx"`
	Signatories []crypto.Signatory `json:"signatories,omitempty"`
	hash        []byte
}

type Signatory struct {
	PublicKey PublicKey `json:"publicKey"`
	Signature Signature `json:"signature"`
}
type Signature struct {
	data signatureData
}

type signatureData struct {
	Signature []byte `json:"signature"`
}


type Tx interface {
	Signers() []TxInput
	Type() Type
	Amount() uint64
	Fee() uint64
	EnsureValid() error
}


type TxInput struct {
	Address  crypto.Address `json:"address"`
	Amount   uint64         `json:"amount"`
	Sequence uint64         `json:"sequence"`
}


	if txs.Txs[0].Tx.Type() == 1 {
		sndTx := txs.Txs[0].Tx
		println(sndTx.Amount())
	}
	println("TX Count: ", txs.Count)
	println("TX Chain ID: ", txs.Txs[0].ChainID)
	println("TX Hash: ", hex.EncodeToString(txs.Txs[0].Hash()))
	println("TX Signatories Public Key: ", txs.Txs[0].Signatories[0].PublicKey.String())
	println("TX Signatories ACC Address: ", txs.Txs[0].Signatories[0].PublicKey.AccountAddress().String())
	println("TX Signatories VAL Address: ", txs.Txs[0].Signatories[0].PublicKey.ValidatorAddress().String())
	println("TX Signatories Signature: ", txs.Txs[0].Signatories[0].Signature.String())
	println("Num Signers: ", len(txs.Txs[0].Tx.Signers()))
	println("TX Signers Address: ", txs.Txs[0].Tx.Signers()[0].Address.String())
	println("TX Signers Amount: ", txs.Txs[0].Tx.Signers()[0].Amount)
	println("TX Signers Sequence: ", txs.Txs[0].Tx.Signers()[0].Sequence)
	println("TX Amount: ", txs.Txs[0].Tx.Amount())
	println("TX Type: ", txs.Txs[0].Tx.Type())
	println("TX Type Striing: ", txs.Txs[0].Tx.Type().String())
	println("TX Fee: ", txs.Txs[0].Tx.Fee())
	println("TX Err: ", txs.Txs[0].Tx.EnsureValid())
*/

//GetTXsCount returns number of TXs
func (g *Burrow) GetTXsCount(height uint64) int {
	return 0

	/*
		client := *g.client
		txs, _ := client.GetBlockTxs(context.Background(), &pb.BlockRequest{Height: height})
		n := int(txs.Count)
		return n
	*/
}

//GetTx returns specified TX
func (g *Burrow) GetTx(height uint64, hash []byte) (*Transaction, error) {
	return nil, nil
	/*
		client := *g.client

		blockRes, getBlockErr := client.GetBlock(context.Background(), &pb.BlockRequest{Height: height})
		if getBlockErr != nil {
			return nil, getBlockErr
		}
		var b Block
		toBlock(blockRes, &b)

		findHash := hex.EncodeToString(hash)
		txRes, getTxErr := client.GetTx(context.Background(), &pb.TxRequest{TxHash: findHash})
		if getTxErr != nil {
			return nil, getTxErr
		}

		var tx Transaction
		toTx(txRes, &tx)
		tx.Time = b.Time
		tx.BlockID = int64(height)

		return &tx, nil
	*/
}

//GetTXs returns all transaction of specific block
func (g *Burrow) GetTXs(height uint64) ([]Transaction, error) {

	url := fmt.Sprintf(g.Domain+"/txs?height=%v", height)
	/*
			response, txsInfoErr := http.Get(url)

			if txsInfoErr != nil {
				println("get txs response error: ", txsInfoErr.Error())
				return nil, txsInfoErr
			}
			response.Close = true
			defer response.Body.Close()

		responseData, _ := ioutil.ReadAll(response.Body)
	*/

	responseData := g.GetReply(url)

	type InputDetails struct {
		Address  string `json:"Address"`
		Amount   uint64 `json:"Amount"`
		Sequence uint64 `json:"Sequence"`
	}

	type OutputDetails struct {
		Address string `json:"Address"`
		Amount  uint64 `json:"Amount"`
	}

	type SendTxPayLoad struct {
		Inputs       []InputDetails           `json:"Inputs"`
		Outputs      []OutputDetails          `json:"Outputs"`
		GasLimit     uint64                   `json:"GasLimit"`
		Fee          uint64                   `json:"Fee"`
		Data         string                   `json:"Data"`
		ContractMeta []map[string]interface{} `json:"ContractMeta"`
	}

	type CallTxPayLoad struct {
		Input        InputDetails             `json:"Input"`
		Address      string                   `json:"Address"`
		GasLimit     uint64                   `json:"GasLimit"`
		Fee          uint64                   `json:"Fee"`
		Data         string                   `json:"Data"`
		ContractMeta []map[string]interface{} `json:"ContractMeta"`
	}

	type TxDetails struct {
		ChainID string                 `json:"ChainID"`
		Type    string                 `json:"Type"`
		Payload map[string]interface{} `json:"Payload"`
	}

	type EnvelopeData struct {
		Signatories []map[string]interface{} `json:"Signatories"`
		Tx          TxDetails                `json:"Tx"`
	}

	type TxData struct {
		Height   uint64      `json:"Height"`
		Hash     string      `json:"Hash"`
		ChainID  string      `json:"ChainID"`
		Payload  interface{} `json:"Payload"`
		Envelope string      `json:"Envelope"`
	}

	type Tx struct {
		Hash string `json:"Hash"`
		Data TxData `json:"Data"`
	}

	type Result struct {
		Count uint64 `json:"Count"`
		Txs   []Tx   `json:"Txs"`
	}

	type BlockTxs struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  Result `json:"result"`
	}

	//var env JsonMsg
	var env BlockTxs
	bytes := []byte(string(responseData))
	if err := json.Unmarshal(bytes, &env); err != nil {
		println("error on get txs from block and unmarshaling response: ", err.Error())
		return nil, err
	}

	if len(env.Result.Txs) <= 0 {
		return nil, fmt.Errorf("no txs exist in this block")
	}

	nTxs := len(env.Result.Txs)
	txs := make([]Transaction, nTxs)

	for i, tx := range env.Result.Txs {

		envStr := tx.Data.Envelope
		/*
			envStr = strings.ReplaceAll(envStr, "\n", "")
			envStr = strings.ReplaceAll(envStr, "\\\"", "\"")
			envStr = strings.ReplaceAll(envStr, "\\\\", "\\")
		*/

		var objEnvelope EnvelopeData
		bytes := []byte(string(envStr))
		if err := json.Unmarshal(bytes, &objEnvelope); err != nil {
			return nil, err
		}

		txs[i].Type = objEnvelope.Tx.Type
		txs[i].BlockID = int64(height)
		txs[i].Hash = tx.Hash

		if objEnvelope.Tx.Type == "CallTx" {

			pl := objEnvelope.Tx.Payload

			input := pl["Input"]
			objInput, ok := input.(map[string]interface{})
			if !ok {
				println("error on converting inputs: ", ok)
				return nil, fmt.Errorf("error on converting inputs")
			}

			Gaslimit := uint64(0)
			if pl["GasLimit"] != nil {
				gaslimitStr := fmt.Sprintf("%v", pl["GasLimit"])
				Gaslimit, _ = strconv.ParseUint(gaslimitStr, 10, 64)
			}

			Fee := uint64(0)
			if pl["Fee"] != nil {
				feeStr := fmt.Sprintf("%v", pl["Fee"])
				Fee, _ = strconv.ParseUint(feeStr, 10, 64)
			}

			data := ""
			if pl["Data"] != nil {
				data = pl["Data"].(string)
			}

			toAddr := ""
			if pl["Address"] != nil {
				toAddr = pl["Address"].(string)
			}

			amount := uint64(0)
			if objInput["Amount"] != nil {
				amountStr := fmt.Sprintf("%v", objInput["Amount"])
				amount, _ = strconv.ParseUint(amountStr, 10, 64)
			}

			seq := uint64(0)
			if objInput["Sequence"] != nil {
				seqStr := fmt.Sprintf("%v", objInput["Sequence"])
				seq, _ = strconv.ParseUint(seqStr, 10, 64)
			}

			txs[i].GasLimit = Gaslimit
			txs[i].Fee = Fee
			txs[i].Data = data
			txs[i].From = objInput["Address"].(string)
			txs[i].To = toAddr
			txs[i].Amount = amount
			txs[i].Sequence = seq

		} else if objEnvelope.Tx.Type == "SendTx" {

			pl := objEnvelope.Tx.Payload

			Inputs := pl["Inputs"]
			objInputs, ok := Inputs.([]interface{}) //map[string]interface{})
			if !ok {
				println("error on converting inputs: ", ok)
				return nil, fmt.Errorf("error on converting inputs")
			}

			fromAddr := ""
			toAddr := ""
			amount := uint64(0)
			seq := uint64(0)

			for i, vi := range objInputs {
				if i == 0 {
					Input, ok := vi.(map[string]interface{})
					if !ok {
						println("error on converting inputs: ", ok)
						return nil, fmt.Errorf("error on converting inputs")
					}
					fromAddr = Input["Address"].(string)

					if Input["Amount"] != nil {
						amountStr := fmt.Sprintf("%v", Input["Amount"])
						amount, _ = strconv.ParseUint(amountStr, 10, 64)
					}

					if Input["Sequence"] != nil {
						seqStr := fmt.Sprintf("%v", Input["Sequence"])
						seq, _ = strconv.ParseUint(seqStr, 10, 64)
					}
				}
			}

			Outputs := pl["Outputs"]
			objOutputs, ok := Outputs.([]interface{})
			if !ok {
				println("error on converting outputs: ", ok)
				return nil, fmt.Errorf("error on converting outputs")
			}
			for i, vi := range objOutputs {
				if i == 0 {
					Output, ok := vi.(map[string]interface{})
					if !ok {
						println("error on converting outputs: ", ok)
						return nil, fmt.Errorf("error on converting outputs")
					}
					toAddr = Output["Address"].(string)
				}
			}

			txs[i].From = fromAddr
			txs[i].To = toAddr
			txs[i].Amount = amount
			txs[i].Sequence = seq

		} else {
			txs[i].From = ""
			txs[i].To = ""
			txs[i].Amount = 0
			txs[i].Sequence = 0
		}

	}

	return txs, nil

}