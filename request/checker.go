package request

import (
	"fmt"
	"github.com/motikan2010/guardian/data"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/motikan2010/guardian/models"
	"github.com/motikan2010/guardian/waf/engine"
)

var staticSuffix = []string{".js", ".css", ".png", ".jpg", ".gif", ".bmp", ".svg", ".ico"}

/*Checker Cheks the requests init*/
type Checker struct {
	ResponseWriter http.ResponseWriter
	Target         *models.Target
	Transaction    *engine.Transaction

	ruleExecutionResult *models.RuleExecutionResult
	startTime           time.Time
}

/*NewRequestChecker Request checker initializer*/
func NewRequestChecker(w http.ResponseWriter, r *http.Request, target *models.Target) *Checker {
	return &Checker{w, target, engine.NewTransaction(r), nil, time.Now()}
}

/*Handle Request checker handler func*/
func (r *Checker) Handle() bool {

	if !r.Target.WAFEnabled || r.Transaction.Request.Method == "GET" && r.IsStaticResource(r.Transaction.Request.URL.Path) {
		return false
	}

	result := r.handleWAFChecker(models.Phase1)

	if result {
		return result
	}

	return r.handleWAFChecker(models.Phase2)
}

// IsStaticResource ...
func (r *Checker) IsStaticResource(url string) bool {
	if strings.Contains(url, "?") {
		return false
	}
	for _, suffix := range staticSuffix {
		if strings.HasSuffix(url, suffix) {
			return true
		}
	}
	return false
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
				r.ruleExecutionResult = &models.RuleExecutionResult{MatchResult: matchResult, Rule: rule}
				break
			} else if !matchResult.IsMatched && !matchResult.DefaultState && !rule.ShouldBlock() {
				matchResult.SetMatch(true)
				r.ruleExecutionResult = &models.RuleExecutionResult{MatchResult: matchResult, Rule: rule}
				break
			}

			//Line 127 and below commented lines are for calculating each rulr exec time
			//passed := helpers.CalcTime(ruleStartTime, time.Now())
			//if passed > 0 {
			//	fmt.Println(rule.Action.ID + " took " + strconv.FormatInt(passed, 10) + " ms.")
			//}
		}

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Minute):
		panic("failed to execute rules.")
	}

	// Detect request
	if r.ruleExecutionResult != nil && r.ruleExecutionResult.MatchResult.IsMatched {

		// Logging
		db := &data.DBHelper{}
		go db.LogMatchResult(r.ruleExecutionResult, r.ruleExecutionResult.Rule.Action.ID, r.Target, r.Transaction.Request.RequestURI, false)

		if r.ruleExecutionResult.Rule.Action.LogAction == models.ActionBlock {
			// Block
			r.ResponseWriter.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(r.ResponseWriter, "Bad Request. %s", url.QueryEscape(r.Transaction.Request.URL.Path))
			return true
		}
	}

	return false
}
