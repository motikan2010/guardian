package engine

import (
	"github.com/motikan2010/guardian/matches"
	"github.com/motikan2010/guardian/waf/bodyprocessor"
)

func init() {
	TransactionMaps.variableMap["REQBODY_PROCESSOR"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			var body string

			switch executer.transaction.RequestBodyProcessor.(type) {
			case *bodyprocessor.JSONBodyProcessor:
				body = "JSON"
				break
			case *bodyprocessor.MultipartProcessor:
				body = "MULTIPART"
				break
			case *bodyprocessor.URLEncodedProcessor:
				body = "URLENCODED"
				break
			case *bodyprocessor.XMLBodyProcessor:
				body = "XML"
				break
			}

			return executer.rule.ExecuteRule(body)
		}}
}
