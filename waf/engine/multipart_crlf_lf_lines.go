package engine

import (
	"github.com/motikan2010/guardian/matches"
)

func init() {
	TransactionMaps.variableMap["MULTIPART_CRLF_LF_LINES"] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			//TODO Not implemented yet
			return matches.NewMatchResult()
		}}
}
