package telegram

import (
	"context"

	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func (r *Runner) resolveStartView(
	ctx context.Context,
	svc *flow.Service,
	spec manager.BotSpec,
	msg *models.Message,
) (flow.ViewModel, error) {
	return svc.Start(ctx, buildStartRequest(spec, msg))
}

func (r *Runner) resolveReplyView(
	ctx context.Context,
	svc *flow.Service,
	spec manager.BotSpec,
	msg *models.Message,
) (flow.ViewModel, bool, error) {
	actionID, ok := svc.ResolveReplyAction(msg.Text)
	if !ok {
		return flow.ViewModel{}, false, nil
	}

	vm, err := svc.HandleAction(
		ctx,
		buildMessageActionRequest(spec, msg, actionID),
	)
	if err != nil {
		return flow.ViewModel{}, true, err
	}

	return vm, true, nil
}

func (r *Runner) resolveCallbackView(
	ctx context.Context,
	svc *flow.Service,
	spec manager.BotSpec,
	cq *models.CallbackQuery,
) (flow.ViewModel, flow.ActionID, bool, error) {
	actionID, ok := decodeActionID(cq.Data)
	if !ok {
		return flow.ViewModel{}, "", false, nil
	}

	vm, err := svc.HandleAction(
		ctx,
		buildCallbackActionRequest(spec, cq, actionID),
	)
	if err != nil {
		return flow.ViewModel{}, actionID, true, err
	}

	return vm, actionID, true, nil
}
