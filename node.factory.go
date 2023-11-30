package knot

import (
	"github.com/hootuu/rock"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/tome/vn"
	"github.com/hootuu/utils/errors"
)

func GetNode(vnID vn.ID, chain kt.Chain) (*Node, *errors.Error) {
	n, err := NewNode(vnID, chain, rock.GetIPFS())
	if err != nil {
		return nil, err
	}
	err = n.doStartup()
	if err != nil {
		return nil, err
	}
	return n, nil
}
