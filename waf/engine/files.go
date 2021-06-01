package engine

import (
	"github.com/motikan2010/guardian/matches"
	"github.com/motikan2010/guardian/waf/bodyprocessor"
)

func init() {
	TransactionMaps.variableMap["FILES"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			matchResult := matches.NewMatchResult()

			switch executer.transaction.RequestBodyProcessor.(type) {
			case *bodyprocessor.MultipartProcessor:

				files := executer.transaction.Request.MultipartForm.File

				for _, headers := range files {
					for _, head := range headers {

						matchResult = executer.rule.ExecuteRule(head.Filename)

						if matchResult.IsMatched {
							return matchResult
						}
					}
				}

				return matchResult
			}

			return matchResult.SetMatch(false)
		}}
}
