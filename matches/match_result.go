package matches

import (
	"time"

	"github.com/motikan2010/guardian/helpers"
)

//MatchResult Match result
type MatchResult struct {
	IsMatched    bool
	DefaultState bool
	StartTime    time.Time
	Elapsed      int64
}

//NewMatchResult Inits match result
func NewMatchResult() *MatchResult {
	return &MatchResult{false, true, time.Now(), 0}
}

//SetMatch ...
func (m *MatchResult) SetMatch(isMatched bool) *MatchResult {
	m.IsMatched = isMatched
	m.DefaultState = false

	m.Elapsed = helpers.CalcTime(m.StartTime, time.Now())

	return m
}
