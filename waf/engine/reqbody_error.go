package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["REQBODY_ERROR"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			if executer.transaction.RequestBodyProcessor.HasBodyError() {
				return executer.rule.ExecuteRule("1")
			}

			return executer.rule.ExecuteRule("0")
		}}
}
