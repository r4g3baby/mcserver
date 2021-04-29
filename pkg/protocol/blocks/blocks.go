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
	blocksByID map[int]map[int]string
)

func init() {
	err := json.Unmarshal(blocksFile, &blocks)
	if err != nil {
		panic(err)
	}

	blocksByID = make(map[int]map[int]string)
	for block, data := range blocks {
		for proto, id := range data {
			if a, ok := blocksByID[id]; ok {
				a[proto] = block
			} else {
				blocksByID[id] = map[int]string{
					proto: block,
				}
			}
		}
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

func GetBlock(id int, protocol protocol.Protocol) string {
	lastProto := 0
	lastBlock := "minecraft:air"
	for proto, block := range blocksByID[id] {
		if proto == int(protocol) {
			return block
		}

		if int(protocol) > proto && proto > lastProto {
			lastProto = proto
			lastBlock = block
		}
	}
	return lastBlock
}
