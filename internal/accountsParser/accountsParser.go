package accountsParser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/util"
)

func parseBalance(accountData types.AccountData) float64 {
	var err error

	for {
		client := util.GetClient()

		req := fasthttp.AcquireRequest()

		req.SetRequestURI(fmt.Sprintf("https://claim.hyperlane.foundation/api/check-eligibility?address=%s",
			accountData.AccountAddress.String()))
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
		req.Header.Set("origin", "https://claim.hyperlane.foundation")
		req.Header.SetReferer("https://claim.hyperlane.foundation/")
		req.Header.SetMethod("GET")

		resp := fasthttp.AcquireResponse()

		if err = client.Do(req, resp); err != nil {
			log.Printf("[%d/%d] | %s | Error When Doing Request When Parsing Balance: %s",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, err)

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		if string(resp.Body()) == "Redirecting...\n" {
			log.Printf("[%d/%d] | %s | Blocked Location",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData)

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		isEligible := gjson.Get(string(resp.Body()), "response.isEligible")
		if !isEligible.Exists() || isEligible.Type == gjson.Null {
			log.Printf("[%d/%d] | %s | Wrong Response When Parsing Balance: %s",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, string(resp.Body()))

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		if !isEligible.Bool() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			return 0
		}

		dropAmount := gjson.Get(string(resp.Body()), "response.eligibilities.0.amount")
		if !dropAmount.Exists() {
			log.Printf("[%d/%d] | %s | Wrong Response When Parsing Balance: %s",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, string(resp.Body()))

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)

		return dropAmount.Float()
	}
}

func ParseAccount(accountData types.AccountData) {
	accountBalance := parseBalance(accountData)

	log.Printf("[%d/%d] | %s | %g $HYPER",
		global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, accountBalance)

	if accountBalance > 0 {
		util.AppendFile("with_balances.txt",
			fmt.Sprintf("%s | %g $KILO\n", accountData.AccountLogData, accountBalance))
	}
}
