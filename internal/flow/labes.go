package flow

const (
	labelsBack = "⬅️ Назад"
)

func backButton() ActionButton {
	return ActionButton{
		ID:    ActionBack,
		Label: labelsBack,
	}
}
