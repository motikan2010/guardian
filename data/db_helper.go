package data

import (
	"fmt"
	"time"

	"github.com/motikan2010/guardian/models"
	"github.com/google/uuid"
)

/*DBHelper The database query helper*/
type DBHelper struct {
}

/*NewDBHelper Inits new db helper*/
func NewDBHelper() *DBHelper {
	return new(DBHelper)
}

//LogMatchResult ...
func (h *DBHelper) LogMatchResult(ruleExecutionResult *models.RuleExecutionResult,
	ruleID string, target *models.Target, requestURI string, forResponse bool) {
	fmt.Printf("--- LogMatchResult ---\n%s\t%s\t%d\t%s\t%s\n",
		uuid.New(),
		time.Now(),
		ruleExecutionResult.MatchResult.Elapsed,
		ruleID,
		requestURI)
}

//LogHTTPRequest ...
func (h *DBHelper) LogHTTPRequest(log *models.HTTPLog) {
	fmt.Printf("--- LogHTTPRequest ---\n%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n",
		log.TargetID,
		log.RequestURI,
		log.StatusCode,
		log.RequestRulesCheckElapsed,
		log.ResponseRulesCheckElapsed,
		log.HTTPElapsed,
		log.RequestSize,
		log.ResponseSize)
}

//LogThrottleRequest ...
func (h *DBHelper) LogThrottleRequest(ipAddress string) {
	fmt.Printf("--- LogThrottleRequest ---\n%s%s%s%d\n",
		uuid.New(),
		time.Now(),
		ipAddress,
		1)
}
