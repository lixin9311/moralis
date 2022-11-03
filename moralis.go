package moralis

import "encoding/json"

/*
{
    "abi": [],
    "block": {
        "number": "",
        "hash": "",
        "timestamp": ""
    },
    "txs": [],
    "txsInternal": [],
    "logs": [],
    "chainId": "",
    "confirmed": true,
    "retries": 0,
    "tag": "",
    "streamId": "",
    "erc20Approvals": [],
    "erc20Transfers": [],
    "nftApprovals": {
        "ERC721": [],
        "ERC1155": []
    },
    "nftTransfers": []
}
*/

type NotifyEvent struct {
	Abi            []interface{}    `json:"abi"`
	Block          Block            `json:"block"`
	Txs            []*Transaction   `json:"txs"`
	TxsInternal    *json.RawMessage `json:"txsInternal"`
	Logs           []*Log           `json:"logs"`
	ChainId        string           `json:"chainId"`
	Confirmed      bool             `json:"confirmed"`
	Retries        int              `json:"retries"`
	Tag            string           `json:"tag"`
	StreamId       string           `json:"streamId"`
	Erc20Approvals *json.RawMessage `json:"erc20Approvals"`
	Erc20Transfers []*Erc20Transfer `json:"erc20Transfers"`
	NftApprovals   *json.RawMessage `json:"nftApprovals"`
	NftTransfers   []NftTransfer    `json:"nftTransfers"`
}

type Block struct {
	Number    int64  `json:"number,string"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp,string"` // timestamp in unix seconds
}

/*
{
	"hash": "0xed781494bdcdc47175d1bc0d870f17d7bbd5815515635e6439a8139e82eba858",
	"gas": "81753",
	"gasPrice": "26000000000",
	"nonce": "17",
	"input": "0xa22cb46500000000000000000000000030f3e83a45fbd610322e3b5e883349da97b2ea030000000000000000000000000000000000000000000000000000000000000001",
	"transactionIndex": "1",
	"fromAddress": "0xd8c12e31fd16e1c9e950760931cc66276e0c2090",
	"toAddress": "0x09ed394ee9bf6c48e728250e55b4fd7f7bb3b417",
	"value": "0",
	"type": "2",
	"v": "0",
	"r": "25714603566502496524590224755204506047701518666282630813488918547350496490761",
	"s": "47295533054968572197327941137220031072488173786145963825697792384055917020645",
	"receiptCumulativeGasUsed": "185394",
	"receiptGasUsed": "54145",
	"receiptContractAddress": null,
	"receiptRoot": null,
	"receiptStatus": "1"
}
*/

type Transaction struct {
	Hash             string `json:"hash"`
	Gas              int64  `json:"gas,string"`
	GasPrice         string `json:"gasPrice"` // big number in string
	Nonce            int64  `json:"nonce,string"`
	Input            string `json:"input"` // transaction payload in hex
	TransactionIndex int64  `json:"transactionIndex,string"`
	FromAddress      string `json:"fromAddress"`
	ToAddress        string `json:"toAddress"`
	Value            string `json:"value"` // big number in string
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
	StreamId         string `json:"streamId"`
	// ReceiptCumulativeGasUsed string `json:"receiptCumulativeGasUsed"`
	// ReceiptGasUsed           string `json:"receiptGasUsed"`
	// ReceiptContractAddress   string `json:"receiptContractAddress"`
	// ReceiptRoot              string `json:"receiptRoot"`
	// ReceiptStatus            string `json:"receiptStatus"`
}

/*
{
	"logIndex": "48",
	"transactionHash": "0xb9730dd1b49061f3b5a6f93e0a66a03be199cad6f21ba5e8747a8087754e3e",
	"transactionIndex": "40",
	"transactionValue": "0",
	"address": "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
	"data": "0x0000000000000000000000000000000000000000204f8a5f22b432605d238000",
	"topic0": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	"topic1": "0x00000000000000000000000068b3f12d6e8d85a8d3dbbc15bba9dc5103b888a4",
	"topic2": "0x0000000000000000000000002faf487a4414fe77e2327f0bf4ae2a264a776ad2",
	"topic3": null
}
*/

type Log struct {
	LogIndex         int64  `json:"logIndex,string"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex int64  `json:"transactionIndex,string"`
	TransactionValue string `json:"transactionValue"`
	Address          string `json:"address"`
	Data             string `json:"data"`
	Topic0           string `json:"topic0"`
	Topic1           string `json:"topic1"`
	Topic2           string `json:"topic2"`
	Topic3           string `json:"topic3"`
}

/*
{
	"from": "0x68b3f12d6e8d85a8d3dbbc15bba9dc5103b888a4",
	"to": "0x2faf487a4414fe77e2327f0bf4ae2a264a776ad2",
	"value": "9999678895548600000000000000",
	"transactionHash": "0xb9730dd1b49061f3b5a6f93e0a66a03be199cad6f21ba5e8747a8087754e3e",
	"transactionIndex": "40",
	"logIndex": "48",
	"contract": "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
	"tokenName": "SHIBA INU",
	"tokenSymbol": "SHIB",
	"tokenDecimals": "18",
	"valueWithDecimals": "9999678895.5486"
}
*/

type Erc20Transfer struct {
	From              string  `json:"from"`
	To                string  `json:"to"`
	Value             string  `json:"value"`
	TransactionHash   string  `json:"transactionHash"`
	TransactionIndex  int64   `json:"transactionIndex,string"`
	LogIndex          int64   `json:"logIndex,string"`
	Contract          string  `json:"contract"`
	TokenName         string  `json:"tokenName"`
	TokenSymbol       string  `json:"tokenSymbol"`
	TokenDecimals     int64   `json:"tokenDecimals,string"`
	ValueWithDecimals float64 `json:"valueWithDecimals,string"`
}

/*
{
	"operator": null,
	"from": "0xd8c12e31fd16e1c9e950760931cc66276e0c2090",
	"to": "0x30f3e83a45fbd610322e3b5e883349da97b2ea03",
	"tokenId": "199",
	"amount": "1",
	"transactionHash": "0x55b9c514af1f5bcb76ba79cac9591268fbc52634e747bdc9da627e56c46ee2fc",
	"logIndex": "13",
	"contract": "0x09ed394ee9bf6c48e728250e55b4fd7f7bb3b417",
	"tokenName": "RunBlox",
	"tokenSymbol": "SHOE",
	"tokenContractType": "ERC721"
}
*/

type NftTransfer struct {
	Operator          *json.RawMessage `json:"operator"` // never seen it in the wild
	From              string           `json:"from"`
	To                string           `json:"to"`
	TokenId           string           `json:"tokenId"`
	Amount            string           `json:"amount"` // always 1 for ERC721, big number for ERC1155
	TransactionHash   string           `json:"transactionHash"`
	LogIndex          int64            `json:"logIndex,string"`
	Contract          string           `json:"contract"`
	TokenName         string           `json:"tokenName"`
	TokenSymbol       string           `json:"tokenSymbol"`
	TokenContractType string           `json:"tokenContractType"` // ERC721 or ERC1155
}
