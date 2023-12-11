package main

import (
	"fmt"
	"github.com/hootuu/knot"
	"github.com/hootuu/rock"
	"github.com/hootuu/tome/nd"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/peer"
)

func init() {

}

func main() {
	peer.Running(func() *errors.Error {
		nd.Init(nd.ID(rock.GetIPFS().Identity.String()), "0xf07b88a2bba771b2b9d141589a8d179cff9ea4de257e55c833c6e7dfbe3deb27")
		_, err := knot.NewRope("testVN", "testChain")
		if err != nil {
			fmt.Println(err)
			return err

		}
		return nil
	}, nil)
}
