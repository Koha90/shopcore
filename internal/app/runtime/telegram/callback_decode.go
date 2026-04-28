package telegram

import (
	"strings"

	"github.com/koha90/shopcore/internal/flow"
)

func decodeCallbackActionID(data string) (flow.ActionID, bool) {
	if !strings.HasPrefix(data, callbackPrefix) {
		return "", false
	}

	raw := strings.TrimPrefix(data, callbackPrefix)
	if raw == "" {
		return "", false
	}

	return flow.ActionID(raw), true
}
