package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["REQUEST_BODY_LENGTH"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {

			bodyBufferLen := len(executer.transaction.RequestBodyProcessor.GetBodyBuffer())

			return executer.rule.ExecuteRule(bodyBufferLen)
		}}
}
