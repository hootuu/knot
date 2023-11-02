package knotx

import (
	"fmt"
)

const (
	MaxBufferSize        = 1024
	Genesis       string = "0"
	HeadRange     int64  = 5
	TailRange     string = "6"

	GenesisPrefix = "G"
	HeadPrefix    = "H"
	TailPrefix    = "T"

	GenesisKey = GenesisPrefix + Genesis
	TailKey    = TailPrefix + TailRange
)

func GetKnotKey(knotIdx int64) string {
	return fmt.Sprintf("%d", knotIdx)
}

func GetGenesisKey() string {
	return GetKnotKey(0)
}

func GetTailKey() string {
	return TailKey
}

//
//func NewGenesisKnot(vnID vn.ID, chain bk.Chain) (*bk.Knot, *errors.Error) {
//	genesisData := bk.NewGenesis(vnID)
//	genesisBID, _, err := rock.Inscribe(genesisData, true)
//	if err != nil {
//		return nil, err
//	}
//	return &bk.Knot{
//		Vn:       vnID,
//		Knot:     bk.GenesisKnotIDX,
//		Previous: bk.GenesisBID,
//		Chain:    chain,
//		Data:     genesisBID,
//	}, nil
//}
