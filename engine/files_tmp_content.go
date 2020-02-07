package engine

import (
	"github.com/asalih/guardian/matches"
)

var FILES_TMP_CONTENT = "FILES_TMP_CONTENT"

func (t *Transaction) loadFilesTmpContent() *Transaction {
	t.variableMap[FILES_TMP_CONTENT] =
		&TransactionData{func(executer *TransactionExecuterModel) *matches.MatchResult {
			//TODO Not implemented yet
			return matches.NewMatchResult(false)
		}}

	return t
}
