package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
)

// saveCatalogImagePhoto downloads the largest Telegram photo and stores it as a
// catalog image filw.
//
// It returns a relative path that can be stored in catalog image_url.
func (r *Runner) saveCatalogImagePhoto(
	ctx context.Context,
	b *tgbot.Bot,
	target flow.CatalogImageInputTarget,
	msg *models.Message,
) (string, error) {
	fileID := telegramPhotoFileToken(msg)
	if fileID == "" {
		return "", fmt.Errorf("telegram photo file id is empty")
	}

	path, err := catalogImageUploadPath(target, time.Now())
	if err != nil {
		return "", err
	}

	file, err := b.GetFile(ctx, &tgbot.GetFileParams{
		FileID: fileID,
	})
	if err != nil {
		return "", fmt.Errorf("get telegram file: %w", err)
	}

	url := b.FileDownloadLink(file)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("build telegram file request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download telegram file: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.log.Error("close telegram file response body failed", "err", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download telegram file: unexpected status %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create catalog image dir: %w", err)
	}

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create catalog image file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			r.log.Error("close catalog image file failed", "err", err)
		}
	}()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return "", fmt.Errorf("write catalog image file: %w", err)
	}

	return path, nil
}

// catalogImageUploadPath builds a relative catalog image path for uploaded admin photos.
//
// The path is stored in catalog image_url and later rendered by Telegram/Web adapters.
// Runtime owns the file system details; flow only receives the final path.
func catalogImageUploadPath(target flow.CatalogImageInputTarget, now time.Time) (string, error) {
	if target.EntityID <= 0 {
		return "", fmt.Errorf("catalog image target entity id is invalid")
	}

	stem := catalogImageFileStem(target)

	switch target.Kind {
	case flow.CatalogImageTargetProduct:
		return fmt.Sprintf(
			"assets/catalog/products/%s-%d.jpg",
			stem,
			now.Unix(),
		), nil

	case flow.CatalogImageTargetVariant:
		return fmt.Sprintf(
			"assets/catalog/variants/%s-%d.jpg",
			stem,
			now.Unix(),
		), nil

	default:
		return "", fmt.Errorf("unknown catalog image target kind %q", target.Kind)
	}
}

func catalogImageFileStem(target flow.CatalogImageInputTarget) string {
	code := sanitizeCatalogImageCode(target.EntityCode)
	if code == "" {
		return fmt.Sprintf("%d", target.EntityID)
	}

	return fmt.Sprintf("%d-%s", target.EntityID, code)
}

func sanitizeCatalogImageCode(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}

	return b.String()
}
