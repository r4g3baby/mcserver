package blocks

import (
	_ "embed"
	"encoding/json"
	"github.com/r4g3baby/mcserver/pkg/protocol"
)

var (
	//go:embed blocks.json
	blocksFile []byte
	blocks     map[string]map[int]int
)

func init() {
	err := json.Unmarshal(blocksFile, &blocks)
	if err != nil {
		panic(err)
	}
}

func GetBlockID(block string, protocol protocol.Protocol) int {
	lastProto := 0
	lastID := 0
	for proto, id := range blocks[block] {
		if proto == int(protocol) {
			return id
		}

		if int(protocol) > proto && proto > lastProto {
			lastProto = proto
			lastID = id
		}
	}
	return lastID
}
