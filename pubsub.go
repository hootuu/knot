package knot

import (
	"fmt"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/tome/vn"
)

type Topic string

func NewTopic(vnID vn.ID, chain kt.Chain) Topic {
	return Topic(fmt.Sprintf("HOTU.%s.%s", vnID.S(), chain.S()))
}

func (t Topic) S() string {
	return string(t)
}
