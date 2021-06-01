package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["ARGS_COMBINED_SIZE"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			sizeOfParams := 0

			queries := executer.transaction.Request.URL.Query()
			for q := range queries {
				if executer.variable.ShouldPassCheck(q) {
					continue
				}

				sizeOfParams += len(queries[q])
			}

			form := executer.transaction.RequestBodyProcessor.GetBody()

			for f := range form {
				if executer.variable.ShouldPassCheck(f) {
					continue
				}

				sizeOfParams += len(form[f])
			}

			return executer.rule.ExecuteRule(sizeOfParams)
		}}
}
