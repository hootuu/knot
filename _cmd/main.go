package main

import (
	"fmt"
	"github.com/hootuu/knot"
	"github.com/hootuu/rock"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/tome/nd"
	"github.com/hootuu/utils/peer"
	"time"
)

func main() {
	go peer.Running()
	time.Sleep(10 * time.Second)
	nd.Init(nd.ID(rock.GetIPFS().Identity.String()), "0xf07b88a2bba771b2b9d141589a8d179cff9ea4de257e55c833c6e7dfbe3deb27")
	rope, err := knot.NewRope("testVN", "testChain")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 1; i++ {
		tpl := &kt.Template{
			Type:      "tpl",
			Version:   kt.Version(i),
			Vn:        "testVN",
			Signature: nil,
		}
		err = kt.InvariableSign(tpl, nd.Here().PRI)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = rope.Submit(tpl)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1000 * time.Second)
}
