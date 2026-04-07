package telegram

import (
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func buildStartRequest(spec manager.BotSpec, msg *models.Message, canAdmin bool) flow.StartRequest {
	return flow.StartRequest{
		BotID:         spec.ID,
		BotName:       spec.Name,
		StartScenario: spec.StartScenario,
		SessionKey:    messageSessionKey(spec.ID, msg),
		CanAdmin:      canAdmin,
	}
}

func buildMessageActionRequest(
	spec manager.BotSpec,
	msg *models.Message,
	actionID flow.ActionID,
	canAdmin bool,
) flow.ActionRequest {
	return flow.ActionRequest{
		BotID:         spec.ID,
		BotName:       spec.Name,
		StartScenario: spec.StartScenario,
		ActionID:      actionID,
		SessionKey:    messageSessionKey(spec.ID, msg),
		CanAdmin:      canAdmin,
	}
}

func buildCallbackActionRequest(
	spec manager.BotSpec,
	cq *models.CallbackQuery,
	actionID flow.ActionID,
	canAdmin bool,
) flow.ActionRequest {
	return flow.ActionRequest{
		BotID:         spec.ID,
		BotName:       spec.Name,
		StartScenario: spec.StartScenario,
		ActionID:      actionID,
		SessionKey:    callbackSessionKey(spec.ID, cq),
		CanAdmin:      canAdmin,
	}
}
