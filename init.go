package knot

import (
	"github.com/hootuu/rock"
	"github.com/hootuu/tome/kt"
	"github.com/hootuu/utils/logger"
)

var gLogger = logger.GetLogger("knot")

func init() {
	rock.RegisterDT([]interface{}{
		kt.Template{},
		kt.Signature{},
	})
}
