package global

import (
	"github.com/valyala/fasthttp"
	"main/pkg/types"
)

var AccountsList []types.AccountData
var Clients []*fasthttp.Client
var TargetProgress = 0
var CurrentProgress = 1
