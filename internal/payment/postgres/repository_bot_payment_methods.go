package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	paymentservice "github.com/koha90/shopcore/internal/payment/service"
)

// ListBotPaymentMethods returns active payment methods enabled for one bot.
func (r *Repository) ListBotPaymentMethods(
	ctx context.Context,
	botID string,
) ([]paymentservice.BotPaymentMethod, error) {
	const op = "payment postgres repository list bot payment methods"

	if err := r.ensureReady(op); err != nil {
		return nil, err
	}

	const q = `
		select
				bpm.id,
				pm.code,
				pm.name,
				pm.kind,
				bpm.display_name,
				bpm.extra_percent_bps,
				bpm.sort_order
		from bot_payment_methods bpm
		join payment_methods pm on pm.id = bpm.payment_method_id
		where bpm.bot_id = $1
				and bpm.is_active = true
				and pm.is_active = true
		order by bpm.sort_order asc, pm.sort_order asc, bpm.id asc
	`

	rows, err := r.pool.Query(ctx, q, botID)
	if err != nil {
		return nil, fmt.Errorf("%s: query bot %q: %w", op, botID, err)
	}

	out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (paymentservice.BotPaymentMethod, error) {
		var item paymentservice.BotPaymentMethod

		if err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.Kind,
			&item.DisplayName,
			&item.ExtraPercentBPS,
			&item.SortOrder,
		); err != nil {
			return paymentservice.BotPaymentMethod{}, err
		}

		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect bot %q: %w", op, botID, err)
	}

	return out, nil
}
