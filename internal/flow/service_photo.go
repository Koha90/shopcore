package flow

import (
	"context"
	"errors"
	"strings"
)

// HandlePhoto resolves a photo message relative to current session state.
//
// It is intentionally separate from HandleText becouse photo input is a
// transport-level media event, not a plain text continuation.
func (s *Service) HandlePhoto(ctx context.Context, req PhotoRequest) (ViewModel, error) {
	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return ViewModel{}, err
	}

	session, ok := s.store.Get(req.SessionKey)
	if !ok {
		session = Session{
			Current:  startScreenForScenario(req.StartScenario),
			History:  nil,
			Pending:  PendingInput{},
			CanAdmin: req.CanAdmin,
		}
	} else {
		session = s.syncSessionAccess(req.SessionKey, session, req.CanAdmin, req.StartScenario)
	}

	if !session.CanAdmin && isAdminPending(session.Pending.Kind) {
		return ViewModel{}, ErrUnknownAction
	}

	if !session.Pending.Active() {
		return s.renderScreen(catalog, session, req.CanAdmin), nil
	}

	switch session.Pending.Kind {
	case PendingInputAdminCustomerPhotoReply:
		fileToken := strings.TrimSpace(req.FileToken)
		chatID, ok := pendingCustomerChatID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending customer chat id is invalid")
		}

		userID, ok := pendingCustomerUserID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending customer user id is invalid")
		}

		if fileToken == "" {
			return buildAdminCustomerPhotoReplyInputView(
				"Фото не найдено. Отправьте изображение ещё раз.",
				chatID,
				userID,
			), nil
		}

		caption := strings.TrimSpace(req.Caption)

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCustomerReplyDone
		s.store.Put(req.SessionKey, session)

		vm := buildAdminCustomerReplyDoneView()
		vm.Effects = []Effect{
			{
				Kind: EffectSendPhoto,
				Target: EffectTarget{
					ChatID: chatID,
					UserID: userID,
				},
				Text: caption,
				Media: &EffectMedia{
					Kind:      EffectMediaPhoto,
					FileToken: fileToken,
				},
			},
		}

		return vm, nil

	default:
		return ViewModel{}, ErrUnknownPendingInput
	}
}
