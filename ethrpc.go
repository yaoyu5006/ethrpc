package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

type ethError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err ethError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *ethError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// EthRPC - Ethereum rpc client
type EthRPC struct {
	url   string
	Debug bool
}

// NewEthRPC create new rpc client with given url
func NewEthRPC(url string) *EthRPC {
	return &EthRPC{url: url}
}

func (rpc *EthRPC) call(method string, target interface{}, params ...interface{}) error {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if rpc.Debug {
		log.Printf("%s\nRequest: %s\nResponse: %s\n", method, body, data)
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	if target == nil {
		return nil
	}

	if err := json.Unmarshal(resp.Result, target); err != nil {
		return err
	}

	return nil
}

// Web3ClientVersion returns the current client version.
func (rpc *EthRPC) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Web3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *EthRPC) Web3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.call("web3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// NetVersion returns the current network protocol version.
func (rpc *EthRPC) NetVersion() (string, error) {
	var version string

	err := rpc.call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (rpc *EthRPC) NetListening() (bool, error) {
	var listening bool

	err := rpc.call("net_listening", &listening)
	return listening, err
}

// NetPeerCount returns number of peers currenly connected to the client.
func (rpc *EthRPC) NetPeerCount() (int, error) {
	var response string
	if err := rpc.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthProtocolVersion returns the current ethereum protocol version.
func (rpc *EthRPC) EthProtocolVersion() (string, error) {
	var protocolVersion string

	err := rpc.call("eth_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

// EthSyncing returns an object with data about the sync status or false.
// func (rpc *EthRPC) EthSyncing() (bool, error) {

// 	err := rpc.call("eth_syncing",)
// 	return syncing, err
// }

// EthCoinbase returns the client coinbase address
func (rpc *EthRPC) EthCoinbase() (string, error) {
	var address string

	err := rpc.call("eth_coinbase", &address)
	return address, err
}

// EthMining returns true if client is actively mining new blocks.
func (rpc *EthRPC) EthMining() (bool, error) {
	var mining bool

	err := rpc.call("eth_mining", &mining)
	return mining, err
}

// EthHashrate returns the number of hashes per second that the node is mining with.
func (rpc *EthRPC) EthHashrate() (int, error) {
	var response string

	if err := rpc.call("eth_hashrate", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGasPrice returns the current price per gas in wei.
func (rpc *EthRPC) EthGasPrice() (big.Int, error) {
	var response string
	if err := rpc.call("eth_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthAccounts returns a list of addresses owned by client.
func (rpc *EthRPC) EthAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("eth_accounts", &accounts)
	return accounts, err
}

// EthBlockNumber returns the number of most recent block.
func (rpc *EthRPC) EthBlockNumber() (int, error) {
	var response string
	if err := rpc.call("eth_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBalance returns the balance of the account of given address in wei.
func (rpc *EthRPC) EthGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := rpc.call("eth_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthGetTransactionCount returns the number of transactions sent from an address.
func (rpc *EthRPC) EthGetTransactionCount(address, block string) (int, error) {
	var response string

	if err := rpc.call("eth_getTransactionCount", &response, address, block); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByHash returns the number of transactions in a block from a block matching the given block hash.
func (rpc *EthRPC) EthGetBlockTransactionCountByHash(hash string) (int, error) {
	var response string

	if err := rpc.call("eth_getBlockTransactionCountByHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (rpc *EthRPC) EthGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	if err := rpc.call("eth_getBlockTransactionCountByNumber", &response, fmt.Sprintf("0x%x", number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetUncleCountByBlockHash returns the number of uncles in a block from a block matching the given block hash.
func (rpc *EthRPC) EthGetUncleCountByBlockHash(hash string) (int, error) {
	var response string

	if err := rpc.call("eth_getUncleCountByBlockHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number.
func (rpc *EthRPC) EthGetUncleCountByBlockNumber(number int) (int, error) {
	var response string

	if err := rpc.call("eth_getUncleCountByBlockNumber", &response, fmt.Sprintf("0x%x", number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetCompilers returns a list of available compilers in the client.
func (rpc *EthRPC) EthGetCompilers() ([]string, error) {
	compilers := []string{}

	err := rpc.call("eth_getCompilers", &compilers)
	return compilers, err
}