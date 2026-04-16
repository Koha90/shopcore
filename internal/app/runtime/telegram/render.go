package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
)

// callbackPrefix is used to version and namespace callback payloads.
//
// Telegram callback_data is limited, so payloads should stay short.
const callbackPrefix = "a:"

// sendView sends a new Telegram message rendered from flow.ViewModel.
func (r *Runner) sendView(
	ctx context.Context,
	b *tgbot.Bot,
	chatID int64,
	vm flow.ViewModel,
) error {
	if hasImage(vm) {
		return r.sendImageView(ctx, b, chatID, vm)
	}

	return r.sendTextView(ctx, b, chatID, vm)
}

func (r *Runner) sendTextView(
	ctx context.Context,
	b *tgbot.Bot,
	chatID int64,
	vm flow.ViewModel,
) error {
	params := &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   vm.Text,
	}

	replyMarkup, err := r.buildReplyMarkup(vm)
	if err != nil {
		return err
	}
	if replyMarkup != nil {
		params.ReplyMarkup = replyMarkup
	}

	if _, err := b.SendMessage(ctx, params); err != nil {
		return fmt.Errorf("send telegram view: %w", err)
	}

	return nil
}

func (r *Runner) sendImageView(
	ctx context.Context,
	b *tgbot.Bot,
	chatID int64,
	vm flow.ViewModel,
) error {
	const op = "send telegram image view"
	if vm.Media == nil {
		return fmt.Errorf("%s: media is nil", op)
	}

	params := &tgbot.SendPhotoParams{
		ChatID:  chatID,
		Photo:   &models.InputFileString{Data: vm.Media.Source},
		Caption: vm.Text,
	}

	replyMarkup, err := r.buildReplyMarkup(vm)
	if err != nil {
		return err
	}

	if inline, ok := replyMarkup.(*models.InlineKeyboardMarkup); ok {
		params.ReplyMarkup = inline
	}

	if _, err := b.SendPhoto(ctx, params); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// editView edits an existing Telegram message rendered from flow.ViewModel.
//
// It is intended primarily for inline-button flows driven by callback queries.
func (r *Runner) editView(
	ctx context.Context,
	b *tgbot.Bot,
	msg *models.Message,
	vm flow.ViewModel,
) error {
	params := &tgbot.EditMessageTextParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
		Text:      vm.Text,
	}

	replyMarkup, err := r.buildReplyMarkup(vm)
	if err != nil {
		return err
	}

	// For edited messages, inline markup is the main supported case.
	// Reply keyboards are typically sent with new message.
	if inline, ok := replyMarkup.(*models.InlineKeyboardMarkup); ok {
		params.ReplyMarkup = inline
	}

	if _, err := b.EditMessageText(ctx, params); err != nil {
		return fmt.Errorf("edit telegram view: %w", err)
	}

	return nil
}

// buildReplyMarkup converts flow.ViewModel keyboard description into Telegram
// reply markup.
func (r *Runner) buildReplyMarkup(vm flow.ViewModel) (models.ReplyMarkup, error) {
	switch {
	case vm.Inline != nil:
		return buildInlineKeyboard(vm.Inline), nil

	case vm.Reply != nil:
		return buildReplyKeyboard(vm.Reply), nil

	case vm.RemoveReply:
		return &models.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		}, nil

	default:
		return nil, nil
	}
}

func buildInlineKeyboard(kb *flow.InlineKeyboardView) *models.InlineKeyboardMarkup {
	if kb == nil || len(kb.Sections) == 0 {
		return nil
	}

	rows := make([][]models.InlineKeyboardButton, 0)

	for _, section := range kb.Sections {
		cols := section.Columns
		if cols <= 0 {
			cols = 1
		}

		current := make([]models.InlineKeyboardButton, 0, cols)

		for _, action := range section.Actions {
			btn := models.InlineKeyboardButton{
				Text:         action.Label,
				CallbackData: encodeActionID(action.ID),
			}

			current = append(current, btn)

			if len(current) == cols {
				rows = append(rows, current)
				current = make([]models.InlineKeyboardButton, 0, cols)
			}
		}

		if len(current) > 0 {
			rows = append(rows, current)
		}
	}

	if len(rows) == 0 {
		return nil
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func buildReplyKeyboard(kb *flow.ReplyKeyboardView) *models.ReplyKeyboardMarkup {
	if kb == nil || len(kb.Rows) == 0 {
		return nil
	}

	rows := make([][]models.KeyboardButton, 0, len(kb.Rows))

	for _, srcRow := range kb.Rows {
		row := make([]models.KeyboardButton, 0, len(srcRow))
		for _, btn := range srcRow {
			row = append(row, models.KeyboardButton{
				Text: btn.Label,
			})
		}
		rows = append(rows, row)
	}

	return &models.ReplyKeyboardMarkup{
		Keyboard:       rows,
		ResizeKeyboard: true,
		IsPersistent:   true,
	}
}

// encodeActionID converts ActionID into compact callback payload.
func encodeActionID(id flow.ActionID) string {
	return callbackPrefix + string(id)
}

// decodeActionID parser callback payload back into ActionID.
func decodeActionID(data string) (flow.ActionID, bool) {
	if !strings.HasPrefix(data, callbackPrefix) {
		return "", false
	}

	raw := strings.TrimPrefix(data, callbackPrefix)
	if raw == "" {
		return "", false
	}

	return flow.ActionID(raw), true
}

func hasImage(vm flow.ViewModel) bool {
	return vm.Media != nil &&
		vm.Media.Kind == flow.MediaKindImage &&
		vm.Media.Source != ""
}

func messageHasImage(msg *models.Message) bool {
	if msg == nil {
		return false
	}

	return len(msg.Photo) > 0
}
