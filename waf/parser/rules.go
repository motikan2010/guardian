package parser

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/motikan2010/guardian/models"
)

type Rule struct {
	RuleId int             `json:"rule_id"`
	ConditionBase64 string `json:"condition_base64"`
	Mode int               `json:"mode"`
}

var xDirectives = []string{"SecAction", "SecArgumentSeparator", "SecAuditEngine", "SecAuditLog", "SecAuditLog2", "SecAuditLogDirMode",
	"SecAuditLogFormat", "SecAuditLogFileMode", "SecAuditLogParts", "SecAuditLogRelevantStatus", "SecAuditLogStorageDir",
	"SecAuditLogType", "SecCacheTransformations", "SecChrootDir", "SecCollectionTimeout", "SecComponentSignature",
	"SecConnEngine", "SecContentInjection", "SecCookieFormat", "SecCookieV0Separator", "SecDataDir", "SecDebugLog", "SecDebugLogLevel",
	"SecDefaultAction", "SecDisableBackendCompression", "SecHashEngine", "SecHashKey", "SecHashParam", "SecHashMethodRx", "SecHashMethodPm",
	"SecGeoLookupDb", "SecGsbLookupDb", "SecGuardianLog", "SecHttpBlKey", "SecInterceptOnError", "SecMarker", "SecPcreMatchLimit",
	"SecPcreMatchLimitRecursion", "SecPdfProtect", "SecPdfProtectMethod", "SecPdfProtectSecret", "SecPdfProtectTimeout", "SecPdfProtectTokenName",
	"SecReadStateLimit", "SecConnReadStateLimit", "SecSensorId", "SecWriteStateLimit", "SecConnWriteStateLimit", "SecRemoteRules",
	"SecRemoteRulesFailAction", "SecRequestBodyAccess", "SecRequestBodyInMemoryLimit", "SecRequestBodyLimit", "SecRequestBodyNoFilesLimit",
	"SecRequestBodyLimitAction", "SecResponseBodyLimit", "SecResponseBodyLimitAction", "SecResponseBodyMimeType", "SecResponseBodyMimeTypesClear",
	"SecResponseBodyAccess", "SecRuleInheritance", "SecRuleEngine", "SecRulePerfTime", "SecRuleRemoveById", "SecRuleRemoveByMsg",
	"SecRuleRemoveByTag", "SecRuleScript", "SecRuleUpdateActionById", "SecRuleUpdateTargetById", "SecRuleUpdateTargetByMsg",
	"SecRuleUpdateTargetByTag", "SecServerSignature", "SecStatusEngine", "SecStreamInBodyInspection", "SecStreamOutBodyInspection", "SecTmpDir",
	"SecUnicodeMapFile", "SecUnicodeCodePage", "SecUploadDir", "SecUploadFileLimit", "SecUploadFileMode", "SecUploadKeepFiles",
	"SecWebAppId", "SecXmlExternalEntity"}

//InitRulesCollection Rules data initializer
func InitRulesCollection() {
	models.RulesCollection = make(map[int][]*models.Rule)

	req, err := http.NewRequest("GET", "http://host.docker.internal:8000/api/rule/get?application_token=" + os.Getenv("APPLICATION_TOKEN"), nil)
	if err != nil {
		log.Fatal(err)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rules []*Rule
	if err := json.Unmarshal(body, &rules); err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}

	for _, rule := range rules {
		InitRulesCollectionFile(rule.RuleId, rule.ConditionBase64, rule.Mode)
	}

}

//InitRulesCollectionFile Rules data initializer
func InitRulesCollectionFile(ruleId int, conditionBase64 string, mode int) {
	fmt.Printf("InitRulesCollectionFile %d %d \n", ruleId, mode)
	dec, err := base64.StdEncoding.DecodeString(conditionBase64)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(dec)))

	var plainTextRules []string
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 1 || strings.HasPrefix(line, "#") {
			continue
		}

		readLine := strings.ReplaceAll(strings.TrimSuffix(strings.TrimSpace(line), "\r"), "\n", " ")

		if len(readLine) <= 1 {
			continue
		}

		plainTextRules = append(plainTextRules, readLine)
	}

	plainTextRulesLen := len(plainTextRules)

	for i := 0; i < plainTextRulesLen; i++ {
		row := plainTextRules[i]
		if strings.HasPrefix(row, "SecRule") {
			var rule *models.Rule
			rule, i = walk(plainTextRules, i, plainTextRulesLen)
			if mode == 0 {
				rule.Action.LogAction = models.ActionBlock
			}
			models.RulesCollection[int(rule.Action.Phase)] = append(models.RulesCollection[int(rule.Action.Phase)], rule)
		}
	}
}

func walk(plainTextRules []string, i int, plainTextRulesLen int) (*models.Rule, int) {
	row := plainTextRules[i]
	var chainRule *models.Rule
	var chainWait bool
	for {
		li := i + 1

		if li >= plainTextRulesLen {
			break
		}

		lr := plainTextRules[li]

		if strings.HasPrefix(lr, "chain") || strings.HasPrefix(lr, "\"chain") {
			chainWait = true
		}

		if strings.HasPrefix(lr, "SecRule") {
			if chainWait {
				chainRule, i = walk(plainTextRules, li, plainTextRulesLen)
			}
			break
		}

		isDirective := false
		for _, dir := range xDirectives {
			if strings.HasPrefix(lr, dir) {
				isDirective = true
				break
			}
		}

		if isDirective {
			break
		}

		i = li
		row += lr
	}

	rule := parseRule(row)

	if chainRule != nil {
		rule.Chain = chainRule
	}

	return rule, i
}

func parseRule(ruleTxt string) *models.Rule {

	variablesReg := regexp.MustCompile(`SecRule\s(.*?)\s`)
	operatorReg := regexp.MustCompile(`(\"@?.*?\")`)

	variablesMatch := variablesReg.FindString(ruleTxt)
	operatorMatch := operatorReg.FindString(ruleTxt)

	if variablesMatch == "" {
		return nil
	}

	variables := parseVariables(variablesMatch)
	operators := parseOperators(operatorMatch)
	action := parseAction(strings.ReplaceAll(strings.ReplaceAll(ruleTxt, operatorMatch, ""), variablesMatch, ""))

	return models.NewRule(variables, operators, action, nil)
}

func parseVariables(variable string) []*models.Variable {
	variable = strings.ReplaceAll(variable, "SecRule ", "")
	varsSplit := strings.Split(variable, "|")
	var dataVariable []*models.Variable

	for _, vars := range varsSplit {
		varsAndFilter := strings.Split(vars, ":")

		if len(varsAndFilter) > 2 {
			//TODO Malformed rule
			continue
		}

		var v *models.Variable

		isLengthCheck := varsAndFilter[0][0] == '&'
		if len(varsAndFilter) > 1 {
			isNotType := varsAndFilter[0][0] == '!'
			varName := strings.Trim(varsAndFilter[0], " ")

			if isNotType || isLengthCheck {
				varName = varName[1:]
			}

			v = &models.Variable{Name: varName, Filter: strings.Split(strings.Trim(varsAndFilter[1], " "), ","), FilterIsNotType: isNotType, LengthCheckForCollection: isLengthCheck}
		} else {
			varName := strings.Trim(varsAndFilter[0], " ")

			if isLengthCheck {
				varName = varName[1:]
			}

			v = &models.Variable{Name: varName, LengthCheckForCollection: isLengthCheck}
		}

		dataVariable = append(dataVariable, v)
	}

	return dataVariable
}

func parseOperators(operator string) *models.Operator {
	isNotOperator := strings.HasPrefix(operator, `"!`)
	var isOperatorSpec bool

	if isNotOperator {
		isOperatorSpec = strings.HasPrefix(operator, `"!@`)
	} else {
		isOperatorSpec = strings.HasPrefix(operator, `"@`)
	}

	parsedOperator := "rx"
	parsedExpression := ""

	if isOperatorSpec {
		operatorReg := regexp.MustCompile(`@(.*?)(\s|\")`)
		opMatch := operatorReg.FindStringSubmatch(operator)

		opr := strings.NewReplacer("\"", "")

		parsedOperator = opr.Replace(opMatch[1])

		r := strings.NewReplacer(parsedOperator, "")
		parsedExpression = r.Replace(operator)

		parsedExpression = strings.Trim(parsedExpression, "\" ")
		parsedExpression = strings.TrimLeft(parsedExpression, "@! ")

	} else {
		parsedExpression = strings.TrimLeft(strings.Trim(operator, "\""), "")
		parsedExpression = strings.TrimLeft(parsedExpression, "@! ")
	}

	return &models.Operator{Func: parsedOperator, Expression: parsedExpression, OperatorIsNotType: isNotOperator}
}

func parseAction(action string) *models.Action {
	idReg := regexp.MustCompile(`id:(.*?),`)
	phaseReg := regexp.MustCompile(`phase:(.*?),`)
	transformReg := regexp.MustCompile(`t:(.*?),`)

	idRegMatch := idReg.FindStringSubmatch(action)
	idRegIdentified := "-1"

	if len(idRegMatch) > 1 {
		idRegIdentified = idRegMatch[1]
	}

	phaseRegMatch := phaseReg.FindStringSubmatch(action)
	phaseRegIdentified := 1

	if len(phaseRegMatch) > 1 {
		phaseRegIdentified, _ = strconv.Atoi(phaseRegMatch[1])
	}

	disrupAct := models.DisruptiveActionBlock

	if strings.Contains(action, models.DisruptiveActionPass.ToString()+",") {
		disrupAct = models.DisruptiveActionPass
	} else if strings.Contains(action, models.DisruptiveActionDrop.ToString()+",") {
		disrupAct = models.DisruptiveActionDrop
	} else if strings.Contains(action, models.DisruptiveActionDeny.ToString()+",") {
		disrupAct = models.DisruptiveActionDeny
	} else if strings.Contains(action, models.DisruptiveActionProxy.ToString()+",") {
		disrupAct = models.DisruptiveActionProxy
	}

	transformMatch := transformReg.FindAllStringSubmatch(action, -1)
	var transforms []string

	if len(transformMatch) > 0 {
		for _, j := range transformMatch {
			if len(j) > 1 {
				transforms = append(transforms, j[1])
			}
		}

	}

	detectAction := models.ActionLog // Default is logging only

	return &models.Action{ID: idRegIdentified, Phase: models.Phase(phaseRegIdentified - 1), Transformations: transforms, DisruptiveAction: disrupAct, LogAction: detectAction}
}
