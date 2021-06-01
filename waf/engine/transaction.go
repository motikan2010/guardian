package engine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/motikan2010/guardian/matches"
	"github.com/motikan2010/guardian/waf/bodyprocessor"

	"github.com/motikan2010/guardian/models"
)

//TransactionMaps global map variable object
var TransactionMaps *TransactionMap = &TransactionMap{make(map[string]*TransactionData)}

//TransactionMap variable map model
type TransactionMap struct {
	variableMap map[string]*TransactionData
}

//Transaction Main model for request response examine
type Transaction struct {
	Request               *http.Request
	Response              *http.Response
	RequestBodyProcessor  bodyprocessor.IBodyProcessor
	ResponseBodyProcessor bodyprocessor.IBodyProcessor

	duration time.Time
	tx       map[string]interface{}
}

//TransactionData Transaction model
type TransactionData struct {
	executer func(*TransactionExecuterModel) *matches.MatchResult
}

//TransactionExecuterModel Executer model
type TransactionExecuterModel struct {
	transaction *Transaction
	rule        *models.Rule
	variable    *models.Variable
}

// NewTransaction Initiates a new request variable object
func NewTransaction(r *http.Request) *Transaction {
	return &Transaction{r, nil, bodyprocessor.NewBodyProcessor(r), nil, time.Now(), make(map[string]interface{})}
}

//Get the data in transaction data
func (tMap *TransactionMap) Get(key string) *TransactionData {
	return tMap.variableMap[key]
}

//Execute Executes transaction for rule
func (t *Transaction) Execute(rule *models.Rule) *matches.MatchResult {

	var matchResult *matches.MatchResult

	for _, variable := range rule.Variables {
		mapData := TransactionMaps.Get(variable.Name)

		if mapData == nil {
			//TODO log unknown Rule
			fmt.Println("Unrecognized variable: " + variable.Name)
			return nil
		}

		executerModel := &TransactionExecuterModel{t, rule, variable}
		matchResult = mapData.executer(executerModel)

		if matchResult.IsMatched {
			if rule.Chain != nil {
				matchResult = t.Execute(rule.Chain)

				if matchResult == nil {
					continue
				}
			}

			if !variable.FilterIsNotType && !rule.Operator.OperatorIsNotType {
				return matchResult
			} else if !matchResult.DefaultState {
				matchResult.SetMatch(false)
			}

		} else if !matchResult.IsMatched && !matchResult.DefaultState && (variable.FilterIsNotType || rule.Operator.OperatorIsNotType) {
			return matchResult.SetMatch(true)
		}
	}

	return matchResult

}
