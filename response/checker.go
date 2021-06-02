package response

import (
	"fmt"
	"github.com/motikan2010/guardian/data"
	"net/http"
	"net/url"
	"time"

	"github.com/motikan2010/guardian/waf/bodyprocessor"
	"github.com/motikan2010/guardian/waf/engine"

	"github.com/motikan2010/guardian/models"
)

//Checker Response checker
type Checker struct {
	ResponseWriter http.ResponseWriter
	Transaction    *engine.Transaction
	Target         *models.Target

	result    *models.RuleExecutionResult
	startTime time.Time
}

/*NewResponseChecker Request checker initializer*/
func NewResponseChecker(w http.ResponseWriter, t *engine.Transaction, resp *http.Response, target *models.Target) *Checker {
	t.Response = resp
	t.ResponseBodyProcessor = bodyprocessor.NewResponseBodyProcessor(resp)

	return &Checker{w, t, target, nil, time.Now()}
}

/*Handle Request checker handler func*/
func (r *Checker) Handle() bool {

	if !r.Target.WAFEnabled {
		return false
	}

	result := r.handleWAFChecker(models.Phase3)

	if result {
		return result
	}

	return r.handleWAFChecker(models.Phase4)
}

func (r *Checker) handleWAFChecker(phase models.Phase) bool {

	done := make(chan bool, 1)

	go func() {

		rulesInPhase := models.RulesCollection[int(phase)]

		for _, rule := range rulesInPhase {

			//ruleStartTime := time.Now()
			matchResult := r.Transaction.Execute(rule)

			if matchResult == nil {
				continue
			}

			if matchResult.IsMatched && rule.ShouldBlock() {
				r.result = &models.RuleExecutionResult{MatchResult: matchResult, Rule: rule}
				break
			} else if !matchResult.IsMatched && !matchResult.DefaultState && !rule.ShouldBlock() {
				matchResult.SetMatch(true)
				r.result = &models.RuleExecutionResult{MatchResult: matchResult, Rule: rule}
				break
			}

			//Line 127 and below commented lines are for calculating each rulr exec time
			//passed := helpers.CalcTime(ruleStartTime, time.Now())
			//if passed > 0 {
			//      fmt.Println(rule.Action.ID + " took " + strconv.FormatInt(passed, 10) + " ms.")
			//}
		}

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Minute):
		panic("failed to execute rules.")
	}

	// Detect response
	if r.result != nil && r.result.MatchResult.IsMatched {

		// Logging
		db := &data.DBHelper{}
		go db.LogMatchResult(r.result, r.result.Rule.Action.ID, r.Target, r.Transaction.Request.RequestURI, false)

		if r.result.Rule.Action.LogAction == models.ActionBlock {
			// Block
			r.ResponseWriter.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(r.ResponseWriter, "Bad Request. %s", url.QueryEscape(r.Transaction.Request.URL.Path))
			return true
		}
	}

	return false
}
