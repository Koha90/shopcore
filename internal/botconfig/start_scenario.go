package botconfig

import "strings"

const (
	StartScenarioReplyWelcome  = "reply_welcome"
	StartScenarioInlineCatalog = "inline_catalog"
)

func isValidStartScenario(v string) bool {
	switch strings.TrimSpace(v) {
	case StartScenarioReplyWelcome, StartScenarioInlineCatalog:
		return true
	default:
		return false
	}
}

func NormilizeStartScenario(v string) string {
	switch strings.TrimSpace(v) {
	case StartScenarioInlineCatalog:
		return StartScenarioInlineCatalog
	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return StartScenarioReplyWelcome
	}
}
