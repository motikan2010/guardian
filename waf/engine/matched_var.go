package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["MATCHED_VAR"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			//TODO Not implemented yet
			return matches.NewMatchResult()
		}}

	TransactionMaps.variableMap["MATCHED_VAR_NAME"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			//TODO Not implemented yet
			return matches.NewMatchResult()
		}}
}
