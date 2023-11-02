package knotx

import (
	"github.com/hootuu/knot/dbx"
	"github.com/hootuu/rock"
	"github.com/hootuu/tome/bk"
	"github.com/hootuu/tome/vn"
	"github.com/hootuu/utils/errors"
	"sync"
	"time"
)

type ChainX struct {
	vn          vn.ID
	chain       bk.Chain
	coll        *dbx.Collection
	lstKnotIdx  bk.KnotIDX
	lstKnotBID  bk.BID
	lstSyncTime time.Time
	queue       chan *Job
	lock        sync.Mutex
}

func (c *ChainX) GetChain() bk.Chain {
	return c.chain
}

func (c *ChainX) GetTailKnotIDX() bk.KnotIDX {
	return c.lstKnotIdx
}

func (c *ChainX) GetTailKnotBID() bk.BID {
	return c.lstKnotBID
}

func NewChainX(vnID vn.ID, chain bk.Chain) (*ChainX, *errors.Error) {
	dbxM, err := doGetChainDBX(vnID)
	if err != nil {
		return nil, err
	}
	collM, err := dbxM.Collection(chain.S())
	if err != nil {
		return nil, err
	}
	invGenesis := bk.NewInvariableGenesis(vnID, chain)
	genesisBID, err := bk.BuildBID(invGenesis)
	if err != nil {
		return nil, err
	}
	return &ChainX{
		vn:          vnID,
		chain:       chain,
		coll:        collM,
		lstKnotIdx:  bk.GenesisKnotIDX,
		lstKnotBID:  genesisBID,
		lstSyncTime: time.Now(),
		queue:       make(chan *Job, MaxBufferSize),
		lock:        sync.Mutex{},
	}, nil
}

func (c *ChainX) SyncChain() *errors.Error {
	c.lock.Lock()
	defer c.lock.Unlock()
	// todo
	return nil
}

func (c *ChainX) Append(tie *bk.Tie) *errors.Error {
	knotM, err := bk.NewKnot(c, tie)
	if err != nil {
		return err
	}
	knotBID, _, err := rock.Inscribe(knotM, true)
	if err != nil {
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lstKnotIdx = knotM.Knot
	c.lstKnotBID = knotBID
	return nil
}

func (c *ChainX) Start() {
	for {
		job := <-c.queue
		err := c.Append(job.Tie)
		if err != nil {

		}
	}
}
