package engine

import (
	"github.com/motikan2010/guardian/helpers"
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["REQUEST_HEADERS_NAMES"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {

			httpData := helpers.GetHeadersNames(executer.transaction.Request.Header)

			matchResult := matches.NewMatchResult()

			if executer.variable.LengthCheckForCollection {
				lenOfHeaders := 0
				for _, key := range httpData {
					if executer.variable.ShouldPassCheck(key) {
						continue
					}

					lenOfHeaders++
				}

				return executer.rule.ExecuteRule(lenOfHeaders)
			}

			for _, key := range httpData {
				if executer.variable.ShouldPassCheck(key) {
					continue
				}
				matchResult = executer.rule.ExecuteRule(key)

				if matchResult.IsMatched {
					return matchResult
				}
			}

			return matchResult
		}}
}
