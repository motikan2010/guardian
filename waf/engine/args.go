package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["ARGS"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			if executer.variable.LengthCheckForCollection {
				return argsLengthHandler(executer, true, true)
			}
			return argsHandler(executer, true, true)
		}}

	TransactionMaps.variableMap["ARGS_GET"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			if executer.variable.LengthCheckForCollection {
				return argsLengthHandler(executer, true, false)
			}
			return argsHandler(executer, true, false)
		}}

	TransactionMaps.variableMap["ARGS_POST"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			if executer.variable.LengthCheckForCollection {
				return argsLengthHandler(executer, false, true)
			}
			return argsHandler(executer, false, true)
		}}
}

func argsHandler(executer *TransactionExecuterModel, executeGet bool, executePost bool) *matches.MatchResult {
	matchResult := matches.NewMatchResult()
	if executeGet {
		queries := executer.transaction.Request.URL.Query()
		for q := range queries {
			if executer.variable.ShouldPassCheck(q) {
				continue
			}

			for _, value := range queries[q] {
				matchResult = executer.rule.ExecuteRule(value)

				if matchResult.IsMatched {
					return matchResult
				}
			}
		}
	}

	if executePost {

		form := executer.transaction.RequestBodyProcessor.GetBody()

		for f := range form {
			if executer.variable.ShouldPassCheck(f) {
				continue
			}
			for _, value := range form[f] {
				matchResult = executer.rule.ExecuteRule(value)

				if matchResult.IsMatched {
					return matchResult
				}
			}
		}
	}

	return matchResult
}

func argsLengthHandler(executer *TransactionExecuterModel, executeGet bool, executePost bool) *matches.MatchResult {

	lengthOfParams := 0
	if executeGet {
		queries := executer.transaction.Request.URL.Query()
		for q := range queries {
			if executer.variable.ShouldPassCheck(q) {
				continue
			}
			lengthOfParams++
		}
	}

	if executePost {

		form := executer.transaction.RequestBodyProcessor.GetBody()

		for f := range form {
			if executer.variable.ShouldPassCheck(f) {
				continue
			}
			lengthOfParams++
		}
	}

	return executer.rule.ExecuteRule(lengthOfParams)
}
