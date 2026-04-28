package flow

import (
	"strconv"
	"strings"
)

const adminCustomerReplyStartPrefix = "admin:message:reply:start:"

// AdminCustomerReplyStartAction builds action that starts admin reply input.
//
// The action carries customer chat/user identifiers. It stays transport-neutral:
// Telegram only transports this action through callback_data.
func AdminCustomerReplyStartAction(chatID, userID int64) ActionID {
	return ActionID(adminCustomerReplyStartPrefix +
		strconv.FormatInt(chatID, 10) +
		":" +
		strconv.FormatInt(userID, 10))
}

// AdminCustomerReplyStartTarget parses customer target from admin reply action.
func AdminCustomerReplyStartTarget(actionID ActionID) (int64, int64, bool) {
	return parseAdminCustomerReplyStartAction(actionID)
}

func parseAdminCustomerReplyStartAction(actionID ActionID) (int64, int64, bool) {
	raw := string(actionID)
	if !strings.HasPrefix(raw, adminCustomerReplyStartPrefix) {
		return 0, 0, false
	}

	rest := strings.TrimPrefix(raw, adminCustomerReplyStartPrefix)
	parts := strings.Split(rest, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}

	chatID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || chatID == 0 {
		return 0, 0, false
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || userID == 0 {
		return 0, 0, false
	}

	return chatID, userID, true
}

func buildAdminCustomerReplyInputView(validation string, chatID, userID int64) ViewModel {
	text := buildAdminTextWithValidation(
		"Ответ пользователю",
		[]string{
			formatAdminFieldLine("User ID", strconv.FormatInt(userID, 10)),
			formatAdminFieldLine("Chat ID", strconv.FormatInt(chatID, 10)),
		},
		validation,
		"Введите ответ сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
	}
}

func buildAdminCustomerReplyDoneView() ViewModel {
	return ViewModel{
		Text: "Ответ пользователю\n\nОтвет отправлен.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func pendingCustomerChatID(p PendingInput) (int64, bool) {
	return pendingInt64(p, PendingValueCustomerChatID)
}

func pendingCustomerUserID(p PendingInput) (int64, bool) {
	return pendingInt64(p, PendingValueCustomerUserID)
}

func pendingInt64(p PendingInput, key string) (int64, bool) {
	raw := strings.TrimSpace(p.Value(key))
	if raw == "" {
		return 0, false
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value == 0 {
		return 0, false
	}

	return value, true
}

func (s *Service) handleAdminCustomerMessageAction(
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool) {
	chatID, userID, ok := parseAdminCustomerReplyStartAction(req.ActionID)
	if !ok {
		return ViewModel{}, session, false
	}

	next := ScreenAdminCustomerReply
	if next != session.Current {
		session.History = append(session.History, session.Current)
		session.Current = next
	}

	session.Pending = PendingInput{
		Kind: PendingInputAdminCustomerReply,
		Payload: PendingInputPayload{
			PendingValueCustomerChatID: strconv.FormatInt(chatID, 10),
			PendingValueCustomerUserID: strconv.FormatInt(userID, 10),
		},
	}

	return buildAdminCustomerReplyInputView("", chatID, userID), session, true
}
