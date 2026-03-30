package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func TestBuildInlineKeyboard(t *testing.T) {
	kb := &flow.InlineKeyboardView{
		Sections: []flow.ActionSection{
			{
				Columns: 2,
				Actions: []flow.ActionButton{
					{ID: flow.ActionEntity1, Label: "Москва"},
					{ID: flow.ActionEntity2, Label: "СПб"},
					{ID: flow.ActionEntity3, Label: "Казань"},
				},
			},
			{
				Columns: 1,
				Actions: []flow.ActionButton{
					{ID: flow.ActionBack, Label: "Назад"},
				},
			},
		},
	}

	got := buildInlineKeyboard(kb)
	require.NotNil(t, got)
	require.Len(t, got.InlineKeyboard, 3)

	require.Len(t, got.InlineKeyboard[0], 2)
	require.Equal(t, "Москва", got.InlineKeyboard[0][0].Text)
	require.Equal(t, encodeActionID(flow.ActionEntity1), got.InlineKeyboard[0][0].CallbackData)
	require.Equal(t, "СПб", got.InlineKeyboard[0][1].Text)

	require.Len(t, got.InlineKeyboard[1], 1)
	require.Equal(t, "Казань", got.InlineKeyboard[1][0].Text)
	require.Equal(t, encodeActionID(flow.ActionEntity3), got.InlineKeyboard[1][0].CallbackData)

	require.Len(t, got.InlineKeyboard[2], 1)
	require.Equal(t, "Назад", got.InlineKeyboard[2][0].Text)
	require.Equal(t, encodeActionID(flow.ActionBack), got.InlineKeyboard[2][0].CallbackData)
}

func TestBuildInlineKeyboard_NilOrEmpty(t *testing.T) {
	require.Nil(t, buildInlineKeyboard(nil))
	require.Nil(t, buildInlineKeyboard(&flow.InlineKeyboardView{}))
}

func TestBuildInlineKeyboard_NormalizesInvalidColumns(t *testing.T) {
	kb := &flow.InlineKeyboardView{
		Sections: []flow.ActionSection{
			{
				Columns: 0,
				Actions: []flow.ActionButton{
					{ID: flow.ActionEntity1, Label: "Москва"},
					{ID: flow.ActionEntity2, Label: "СПб"},
				},
			},
		},
	}

	got := buildInlineKeyboard(kb)
	require.NotNil(t, got)
	require.Len(t, got.InlineKeyboard, 2)
	require.Len(t, got.InlineKeyboard[0], 1)
	require.Len(t, got.InlineKeyboard[1], 1)
}

func TestBuildReplyKeyboard(t *testing.T) {
	kb := &flow.ReplyKeyboardView{
		Rows: [][]flow.ReplyButton{
			{
				{ID: flow.ActionCatalogStart, Label: "♻️ Каталог"},
				{ID: flow.ActionCabinetOpen, Label: "⚙️ Мой кабинет"},
			},
			{
				{ID: flow.ActionSupportOpen, Label: "🤷‍♂️ Поддержка"},
			},
		},
	}

	got := buildReplyKeyboard(kb)
	require.NotNil(t, got)
	require.True(t, got.ResizeKeyboard)
	require.True(t, got.IsPersistent)

	require.Len(t, got.Keyboard, 2)
	require.Equal(t, "♻️ Каталог", got.Keyboard[0][0].Text)
	require.Equal(t, "⚙️ Мой кабинет", got.Keyboard[0][1].Text)
	require.Equal(t, "🤷‍♂️ Поддержка", got.Keyboard[1][0].Text)
}

func TestBuildReplyKeyboard_NilOrEmpty(t *testing.T) {
	require.Nil(t, buildReplyKeyboard(nil))
	require.Nil(t, buildReplyKeyboard(&flow.ReplyKeyboardView{}))
}

func TestBuildReplyMarkup(t *testing.T) {
	r := &Runner{}

	t.Run("inline has priority", func(t *testing.T) {
		vm := flow.ViewModel{
			Inline: &flow.InlineKeyboardView{
				Sections: []flow.ActionSection{
					{
						Columns: 1,
						Actions: []flow.ActionButton{
							{ID: flow.ActionEntity1, Label: "Москва"},
						},
					},
				},
			},
			Reply: &flow.ReplyKeyboardView{
				Rows: [][]flow.ReplyButton{
					{
						{ID: flow.ActionCatalogStart, Label: "♻️ Каталог"},
					},
				},
			},
			RemoveReply: true,
		}

		got, err := r.buildReplyMarkup(vm)
		require.NoError(t, err)

		inline, ok := got.(*models.InlineKeyboardMarkup)
		require.True(t, ok)
		require.Len(t, inline.InlineKeyboard, 1)
	})

	t.Run("reply keyboard", func(t *testing.T) {
		vm := flow.ViewModel{
			Reply: &flow.ReplyKeyboardView{
				Rows: [][]flow.ReplyButton{
					{
						{ID: flow.ActionCatalogStart, Label: "♻️ Каталог"},
					},
				},
			},
		}

		got, err := r.buildReplyMarkup(vm)
		require.NoError(t, err)

		reply, ok := got.(*models.ReplyKeyboardMarkup)
		require.True(t, ok)
		require.Len(t, reply.Keyboard, 1)
	})

	t.Run("remove reply", func(t *testing.T) {
		vm := flow.ViewModel{
			RemoveReply: true,
		}

		got, err := r.buildReplyMarkup(vm)
		require.NoError(t, err)

		remove, ok := got.(*models.ReplyKeyboardRemove)
		require.True(t, ok)
		require.True(t, remove.RemoveKeyboard)
	})

	t.Run("no markup", func(t *testing.T) {
		got, err := r.buildReplyMarkup(flow.ViewModel{})
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func TestEncodeDecodeActionID(t *testing.T) {
	encoded := encodeActionID(flow.ActionEntity1)
	require.Equal(t, callbackPrefix+string(flow.ActionEntity1), encoded)

	decoded, ok := decodeActionID(encoded)
	require.True(t, ok)
	require.Equal(t, flow.ActionEntity1, decoded)
}

func TestDecodeActionID_Invalid(t *testing.T) {
	tests := []string{
		"",
		"x:entity:1",
		callbackPrefix,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			got, ok := decodeActionID(tt)
			require.False(t, ok)
			require.Empty(t, got)
		})
	}
}
