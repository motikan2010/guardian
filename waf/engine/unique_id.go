package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["UNIQUE_ID"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			//Not implemented yet - might not needed.
			return matches.NewMatchResult().SetMatch(true)
		}}
}
