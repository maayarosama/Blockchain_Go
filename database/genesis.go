package database

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
)

var genesisJson = `
{
  "genesis_time": "2022-01-08T00:00:00.000000000Z",
  "chain_id": "Blockchain_Go",
  "balances": {
    "0x73e4B4d8A69b8820744988A0f7D098cAa6066352": 1000000
  }
}`

type Genesis struct {
	Balances map[common.Address]uint `json:"balances"`
}

func loadGenesis(path string) (Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var loadedGenesis Genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return Genesis{}, err
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(path string, genesis []byte) error {
	return ioutil.WriteFile(path, genesis, 0644)
}
