package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["REQUEST_BODY"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			bodyBuffer := executer.transaction.RequestBodyProcessor.GetBodyBuffer()

			if bodyBuffer == nil {
				return matches.NewMatchResult()
			}

			return executer.rule.ExecuteRule(string(bodyBuffer))
		}}
}
