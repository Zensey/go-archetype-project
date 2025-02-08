package protocol

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexIntUnmarshal(t *testing.T) {
	test := struct {
		ID hexInt `json:"id"`
	}{}

	data := []byte(`{"id": "0xdeadbeef"}`)
	err := json.Unmarshal(data, &test)

	require.Nil(t, err)
	require.Equal(t, hexInt(3735928559), test.ID)
}

func TestTrxUnmarshal(t *testing.T) {
	data := []byte(`
            {
                "blockHash": "0xcf1c6419ac387395c4c7e397e5d42cea203ab7ddc7c38ffc14281ddebb9e12dc",
                "blockNumber": "0x149a9f3",
                "from": "0x36a454aef52938c8637cd4689b2980c1cfd43389",
                "gas": "0x35f50",
                "gasPrice": "0x7531c239",
                "maxPriorityFeePerGas": "0x0",
                "maxFeePerGas": "0xafcaa355",
                "hash": "0x850c48096915adfc2453a64a8f444a61390a783647a7f44f8590b04e89ffdac0",
                "input": "0x78e111f60000000000000000000000009e2bc71f52ee2eee7a3cc6d5a183ddd2b57a816b000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c42f1c6b500000000000000000000000000000000000000000000000000000000439d4aeea000000000000000000000000000000000012027283e67f60000000000000000000000000000000000000000000000000000025c886aa8616cbd409ae10d2fb6e0000000000000000000000000000000000000000004d565ad19f91c0000000000000000000000000000000000000000000000000000000000000000067830e03ff8000000000000000000000000000000000000000000000000000000000fd5f00000000000000000000000000000000000000000000000000000000",
                "nonce": "0x41087",
                "to": "0xa69babef1ca67a37ffaf7a485dfff3382056e78c",
                "transactionIndex": "0x0",
                "value": "0x7b8700",
                "type": "0x2",
                "accessList": [
                    {
                        "address": "0xe3fe800b0de664bf0bca8ad58ecbc73b259047b0",
                        "storageKeys": [
                            "0x0000000000000000000000000000000000000000000000000000000000000009",
                            "0x000000000000000000000000000000000000000000000000000000000000000a",
                            "0x0000000000000000000000000000000000000000000000000000000000000000",
                            "0x0000000000000000000000000000000000000000000000000000000000000004",
                            "0x0000000000000000000000000000000000000000000000000000000000000001",
                            "0x7390b40298862f44ea9311d09f2ed00d1d7dc95d3769e730eead11100ea80a20"
                        ]
                    },
                    {
                        "address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                        "storageKeys": [
                            "0x1ab7ebbbf784c32d3a532fae1f41128509a0b64fa588589ad051677a92a9cac4",
                            "0x75245230289a9f0bf73a6c59aef6651b98b3833a62a3c0bd9ab6b0dec8ed4d8f"
                        ]
                    },
                    {
                        "address": "0x595832f8fc6bf59c85c527fec3740a1b7a361269",
                        "storageKeys": [
                            "0x99713ceb4322a7b2d063a2b1e90a212070b8c507ea9c7afebed78f66997ae15e",
                            "0xe1a4184162dfeb463ec20be8837b5c68f27a2b1737f3c72b04526fe65b9578a3"
                        ]
                    },
                    {
                        "address": "0x9e2bc71f52ee2eee7a3cc6d5a183ddd2b57a816b",
                        "storageKeys": []
                    }
                ],
                "chainId": "0x1",
                "v": "0x1",
                "yParity": "0x1",
                "r": "0xea8887d039efe2b5b31b4c7730e3680cd6eaefd1245c6790255a59a80e8188ca",
                "s": "0x1ab1f3b5ce357cfb9bea1e02d0a01a08d019638739ad641193100f0dc30cde00"
            }	
	`)

	test := Transaction{}
	err := test.Parse(data)

	//err := json.Unmarshal(data, &test)
	// fmt.Println((*big.Int)(&test.GasPrice))

	require.Nil(t, err)
	require.Equal(t, (*big.NewInt(1966195257)), (big.Int)(test.GasPrice))
}

func TestResponseUnmarshal(t *testing.T) {

	data, err := os.ReadFile("./assets/test-data-a.json")
	require.Nil(t, err)

	resp := EthResponse{}
	err = json.Unmarshal(data, &resp)
	require.Nil(t, err)

	block := BlockWithTransactions{}
	err = json.Unmarshal(resp.Result, &block)

    require.Nil(t, err)
	require.Equal(t, len(block.Transactions), 588)
}
