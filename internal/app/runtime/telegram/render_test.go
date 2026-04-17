package telegram

import (
	"os"
	"path/filepath"
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
					{ID: flow.ActionID("catalog:select:city:moscow"), Label: "Москва"},
					{ID: flow.ActionID("catalog:select:city:spb"), Label: "СПб"},
					{ID: flow.ActionID("catalog:select:city:kazan"), Label: "Казань"},
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
	require.Equal(t, encodeActionID(flow.ActionID("catalog:select:city:moscow")), got.InlineKeyboard[0][0].CallbackData)
	require.Equal(t, "СПб", got.InlineKeyboard[0][1].Text)

	require.Len(t, got.InlineKeyboard[1], 1)
	require.Equal(t, "Казань", got.InlineKeyboard[1][0].Text)
	require.Equal(t, encodeActionID(flow.ActionID("catalog:select:city:kazan")), got.InlineKeyboard[1][0].CallbackData)

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
					{ID: flow.ActionID("catalog:select:city:moscow"), Label: "Москва"},
					{ID: flow.ActionID("catalog:select:city:spb"), Label: "СПб"},
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
							{ID: flow.ActionID("catalog:select:city:moscow"), Label: "Москва"},
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
	action := flow.ActionID("catalog:select:city:moscow")

	encoded := encodeActionID(action)
	require.Equal(t, callbackPrefix+string(action), encoded)

	decoded, ok := decodeActionID(encoded)
	require.True(t, ok)
	require.Equal(t, action, decoded)
}

func TestDecodeActionID_Invalid(t *testing.T) {
	tests := []string{
		"",
		"x:city:moscow",
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

func TestHasImage(t *testing.T) {
	t.Parallel()

	require.False(t, hasImage(flow.ViewModel{}))
	require.False(t, hasImage(flow.ViewModel{
		Media: &flow.MediaView{Kind: flow.MediaKindImage},
	}))
	require.True(t, hasImage(flow.ViewModel{
		Media: &flow.MediaView{
			Kind:   flow.MediaKindImage,
			Source: "assets/catalog/variants/classic.png",
		},
	}))
}

func TestMessageHasImage(t *testing.T) {
	t.Parallel()

	require.False(t, messageHasImage(nil))
	require.False(t, messageHasImage(&models.Message{}))
	require.True(t, messageHasImage(&models.Message{
		Photo: []models.PhotoSize{
			{FileID: "photo-1"},
		},
	}))
}

func TestIsRemoteMediaSource(t *testing.T) {
	t.Parallel()

	require.True(t, isRemoteMediaSource("https://example.com/a.png"))
	require.True(t, isRemoteMediaSource("http://example.com/a.png"))
	require.False(t, isRemoteMediaSource("assets/catalog/a.png"))
	require.False(t, isRemoteMediaSource(""))
}

func TestBuildTelegramPhotoInput_Remote(t *testing.T) {
	t.Parallel()

	got, err := buildTelegramPhotoInput("https://example.com/a.png")
	require.NoError(t, err)

	file, ok := got.(*models.InputFileString)
	require.True(t, ok)
	require.Equal(t, "https://example.com/a.png", file.Data)
}

func TestBuildTelegramPhotoInput_Local(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "classic.png")
	require.NoError(t, os.WriteFile(path, []byte("png-data"), 0o644))

	got, err := buildTelegramPhotoInput(path)
	require.NoError(t, err)

	file, ok := got.(*models.InputFileUpload)
	require.True(t, ok)
	require.Equal(t, "classic.png", file.Filename)
	require.NotNil(t, file.Data)
}

func TestBuildTelegramPhotoInput_Empty(t *testing.T) {
	t.Parallel()

	got, err := buildTelegramPhotoInput("")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestBuildTelegramPhotoInput_MissingFile(t *testing.T) {
	t.Parallel()

	got, err := buildTelegramPhotoInput("/no/such/file.png")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestBuildTelegramInputMediaPhoto_Remote(t *testing.T) {
	t.Parallel()

	got, err := buildTelegramInputMediaPhoto("https://example.com/a.png", "caption")
	require.NoError(t, err)

	media, ok := got.(*models.InputMediaPhoto)
	require.True(t, ok)
	require.Equal(t, "https://example.com/a.png", media.Media)
	require.Equal(t, "caption", media.Caption)
}

func TestBuildTelegramInputMediaPhoto_Local(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "classic.png")
	require.NoError(t, os.WriteFile(path, []byte("png-data"), 0o644))

	got, err := buildTelegramInputMediaPhoto(path, "caption")
	require.NoError(t, err)

	media, ok := got.(*models.InputMediaPhoto)
	require.True(t, ok)
	require.Equal(t, "attach://classic.png", media.Media)
	require.Equal(t, "caption", media.Caption)
	require.NotNil(t, media.MediaAttachment)
}

func TestBuildTelegramInputMediaPhoto_Empty(t *testing.T) {
	t.Parallel()

	got, err := buildTelegramInputMediaPhoto("", "caption")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestClassifyRenderTransition(t *testing.T) {
	t.Parallel()

	textMsg := &models.Message{}
	imageMsg := &models.Message{
		Photo: []models.PhotoSize{{FileID: "p1"}},
	}

	textVM := flow.ViewModel{}
	imageVM := flow.ViewModel{
		Media: &flow.MediaView{
			Kind:   flow.MediaKindImage,
			Source: "assets/catalog/variants/classic.png",
		},
	}

	require.Equal(t, renderTransitionEditText, classifyRenderTransition(textMsg, textVM))
	require.Equal(t, renderTransitionEditImage, classifyRenderTransition(textMsg, imageVM))
	require.Equal(t, renderTransitionEditImage, classifyRenderTransition(imageMsg, imageVM))
	require.Equal(t, renderTransitionReplaceImageWithText, classifyRenderTransition(imageMsg, textVM))
}
