package network

import "blockchain_from_scratch/utils"

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{address, kind, items}
	payload := utils.GobEncode(inventory)
	req := append(commandToBytes("inv"), payload...)

	sendData(address, req)
}
