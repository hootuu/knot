package knot

import (
	"context"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/tome/nd"
	"github.com/hootuu/tome/vn"
	"github.com/hootuu/twelve"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/sys"
	"github.com/ipfs/kubo/core"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.uber.org/zap"
	"time"
)

type Node struct {
	Topic     Topic
	listener  twelve.ITwelveListener
	ipfsNode  *core.IpfsNode
	ipfsTopic *pubsub.Topic
	ipfsSub   *pubsub.Subscription
	ctx       context.Context
	cancelFuc context.CancelFunc
}

func NewNode(vnID vn.ID, chain kt.Chain, ipfsNode *core.IpfsNode) (*Node, *errors.Error) {
	n := &Node{
		Topic:     NewTopic(vnID, chain),
		listener:  nil,
		ipfsNode:  ipfsNode,
		ipfsTopic: nil,
		ipfsSub:   nil,
	}
	n.ctx, n.cancelFuc = context.WithCancel(context.Background())
	err := n.doStartup()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Node) Close() {
	n.cancelFuc()
}

func (n *Node) Register(listener twelve.ITwelveListener) {
	n.listener = listener
}

func (n *Node) Notify(letter *twelve.Letter) *errors.Error {
	letterData, err := letter.ToBytes()
	if err != nil {
		gLogger.Error("letter.ToBytes() error", zap.Error(err), zap.Any("letter", letter))
		return err
	}
	nErr := n.ipfsTopic.Publish(context.Background(), letterData)
	if nErr != nil {
		gLogger.Error("n.ipfsTopic.Publish(context.Background(), msgData) error", zap.Error(nErr))
		return errors.Sys("topic.Publish error: " + nErr.Error())
	}
	return nil
}

func (n *Node) Node() *nd.Node {
	return nd.Here()
}

func (n *Node) doWaiting() {
	for {
		msg, nErr := n.ipfsSub.Next(n.ctx)
		if nErr != nil {
			gLogger.Error("n.ipfsSub.Next(n.ctx) error", zap.Error(nErr))
			time.Sleep(1 * time.Second)
			continue
		}
		if msg.Local {
			continue
		}
		n.doDealMessage(msg)
	}
}

func (n *Node) doDealMessage(msg *pubsub.Message) {
	msgBytes := msg.GetData()
	letter, err := twelve.LetterOfBytes(msgBytes)
	if err != nil {
		gLogger.Error("twelve.LetterOfBytes error", zap.Error(err), zap.Any("msg", msg))
		return
	}
	err = n.listener.On(letter)
	if err != nil {
		gLogger.Error("n.listener.On(letter) error", zap.Error(err), zap.Any("msg", msg))
		return
	}
}

func (n *Node) doStartup() *errors.Error {
	if n.ipfsTopic != nil {
		return nil //todo
	}
	sys.Info("Join Topic: ", n.Topic)
	var nErr error
	n.ipfsTopic, nErr = n.ipfsNode.PubSub.Join(n.Topic.S())
	if nErr != nil {
		gLogger.Error("n.ipfsNode.PubSub.Join(n.Topic.S()) error", zap.Error(nErr),
			zap.String("topic", n.Topic.S()))
		return errors.Sys("join topic error: " + nErr.Error())
	}
	n.ipfsSub, nErr = n.ipfsTopic.Subscribe()
	if nErr != nil {
		gLogger.Error("n.ipfsTopic.Subscribe() error", zap.Error(nErr),
			zap.String("topic", n.Topic.S()))
		return errors.Sys("topic.Subscribe error: " + nErr.Error())
	}
	go n.doWaiting()
	return nil
}
