package engine

import (
	"github.com/motikan2010/guardian/matches"
	"github.com/motikan2010/guardian/waf/bodyprocessor"
)

func init() {
	transData := &TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
		matchResult := matches.NewMatchResult()

		switch executer.transaction.RequestBodyProcessor.(type) {
		case *bodyprocessor.MultipartProcessor:

			files := executer.transaction.Request.MultipartForm.File
			for _, headers := range files {
				for _, head := range headers {
					matchResult = executer.rule.ExecuteRule(head.Filename)

					if matchResult.IsMatched {
						return matchResult.SetMatch(true)
					}
				}
			}

			return matchResult
		}

		return matchResult.SetMatch(false)
	}}

	TransactionMaps.variableMap["FILES_NAMES"] = transData
	TransactionMaps.variableMap["MULTIPART_FILENAME"] = transData

}
