package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["IP"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {

			return matches.NewMatchResult()

		}}
}
