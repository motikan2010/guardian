package engine

import (
	"time"

	"github.com/motikan2010/guardian/helpers"
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["DURATION"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {

			duration := helpers.CalcTime(executer.transaction.duration, time.Now())

			return executer.rule.ExecuteRule(duration)
		}}
}
