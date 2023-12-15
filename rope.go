package knot

import (
	"github.com/hootuu/tome/bk/bid"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/tome/nd"
	"github.com/hootuu/tome/vn"
	"github.com/hootuu/twelve"
	"github.com/hootuu/utils/errors"
)

type Rope struct {
	vn     vn.ID
	chain  kt.Chain
	height kt.Height
	tail   kt.KID
	node   *Node
	tw     *twelve.Twelve
}

func NewRope(vnID vn.ID, chain kt.Chain) (*Rope, *errors.Error) {
	r := &Rope{
		vn:     vnID,
		chain:  chain,
		height: 0,
		tail:   "",
		node:   nil,
		tw:     nil,
	}
	var err *errors.Error
	r.node, err = GetNode(vnID, chain)
	if err != nil {
		return nil, err
	}
	r.tw, err = twelve.NewTwelve(vnID, kt.Chain(chain), r.node, &twelve.Option{Expect: 2})
	if err != nil {
		return nil, err
	}
	r.tw.Start()
	return r, nil
}

func (r *Rope) Submit(inv kt.Invariable) *errors.Error {
	//stone.Inscribe(inv, true) TODO
	dataBID, err := bid.Build(inv)
	if err != nil {
		return err
	}
	letter := twelve.NewLetter(inv.GetVN(), r.chain, dataBID, twelve.UnLock, twelve.RequestArrow, r.node.Node().ID)
	err = letter.Sign(nd.Here().PRI)
	if err != nil {
		return err
	}
	err = r.node.Notify(letter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rope) GetChain() kt.Chain {
	return r.chain
}

func (r *Rope) GetHeight() kt.Height {
	return r.height
}

func (r *Rope) GetTail() kt.KID {
	return r.tail
}
