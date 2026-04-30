package service

// PaymentKind identifies method behavior class.
type PaymentKind string

const (
	PaymentKindBankCard    PaymentKind = "bank_card"
	PaymentKindSBP         PaymentKind = "sbp"
	PaymentKindMobilePhone PaymentKind = "mobile_phone"
	PaymentKindBTC         PaymentKind = "btc"
	PaymentKindETH         PaymentKind = "eth"
	PaymentKindCash        PaymentKind = "cash"
	PaymentKindManual      PaymentKind = "manual"
)

// BotPaymentMethod describes one payment method enabled for a bot.
type BotPaymentMethod struct {
	ID              int
	Code            string
	Name            string
	Kind            PaymentKind
	DisplayName     string
	ExtraPercentBPS int
	SortOrder       int
}

// Label returns display name preferred for customer-facing screens.
func (m BotPaymentMethod) Label() string {
	if m.DisplayName != "" {
		return m.DisplayName
	}

	return m.Name
}
