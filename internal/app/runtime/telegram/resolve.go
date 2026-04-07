package telegram

import (
	"context"
	"errors"

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
	if msg == nil {
		return flow.ViewModel{}, errors.New("telegram message is nil")
	}
	if msg.From == nil {
		return flow.ViewModel{}, errors.New("telegram message sender is nil")
	}

	canAdmin := r.canAdminTelegram(spec, msg.From.ID)

	return svc.Start(ctx, buildStartRequest(spec, msg, canAdmin))
}

func (r *Runner) resolveReplyView(
	ctx context.Context,
	svc *flow.Service,
	spec manager.BotSpec,
	msg *models.Message,
) (flow.ViewModel, bool, error) {
	if msg == nil {
		return flow.ViewModel{}, false, errors.New("telegram message is nil")
	}
	if msg.From == nil {
		return flow.ViewModel{}, false, errors.New("telegram message sender is nil")
	}

	actionID, ok := svc.ResolveReplyAction(msg.Text)
	if !ok {
		return flow.ViewModel{}, false, nil
	}

	canAdmin := r.canAdminTelegram(spec, msg.From.ID)

	vm, err := svc.HandleAction(
		ctx,
		buildMessageActionRequest(spec, msg, actionID, canAdmin),
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
	if cq == nil {
		return flow.ViewModel{}, "", false, errors.New("telegram callback query is nil")
	}

	actionID, ok := decodeActionID(cq.Data)
	if !ok {
		return flow.ViewModel{}, "", false, nil
	}

	canAdmin := r.canAdminTelegram(spec, cq.From.ID)

	vm, err := svc.HandleAction(
		ctx,
		buildCallbackActionRequest(spec, cq, actionID, canAdmin),
	)
	if err != nil {
		return flow.ViewModel{}, actionID, true, err
	}

	return vm, actionID, true, nil
}

func messageUserID(msg *models.Message) (int64, bool) {
	if msg == nil || msg.From == nil {
		return 0, false
	}

	return msg.From.ID, true
}

func callbackUserID(cq *models.CallbackQuery) (int64, bool) {
	if cq == nil {
		return 0, false
	}

	return cq.From.ID, true
}
