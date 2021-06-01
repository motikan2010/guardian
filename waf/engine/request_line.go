package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["REQUEST_LINE"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {

			uri := executer.transaction.Request.Host

			if executer.transaction.Request.TLS == nil {
				uri = "http://" + uri
			} else {
				uri = "https://" + uri
			}

			if executer.transaction.Request.RequestURI != "" {
				uri += executer.transaction.Request.RequestURI
			}

			line := executer.transaction.Request.Method + " " + uri
			return executer.rule.ExecuteRule(line)
		}}
}
