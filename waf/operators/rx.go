package operators

import (
	"strconv"

	"github.com/motikan2010/guardian/helpers"
)

func init() {
	OperatorMaps.funcMap["rx"] = func(expression interface{}, variableData interface{}) bool {

		switch variableData.(type) {
		case string:
			isMatch, _ := helpers.IsMatch(expression.(string), variableData.(string))
			return isMatch
		case int:
			isMatch, _ := helpers.IsMatch(expression.(string), strconv.Itoa(variableData.(int)))
			return isMatch
		}

		return false
	}
}
