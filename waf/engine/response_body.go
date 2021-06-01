package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["RESPONSE_BODY"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			bodyBuffer := executer.transaction.ResponseBodyProcessor.GetBodyBuffer()

			if bodyBuffer == nil {
				return matches.NewMatchResult()
			}

			return executer.rule.ExecuteRule(string(bodyBuffer))
		}}
}
