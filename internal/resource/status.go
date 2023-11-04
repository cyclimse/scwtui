package resource

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusActive  Status = "active"
	StatusError   Status = "error"
	StatusDeleted Status = "deleted"
	StatusReady   Status = "ready"
	StatusPending Status = "pending"
)

func (s *Status) Emoji() rune {
	if s == nil {
		return ' '
	}

	switch *s {
	case StatusActive, StatusReady:
		return '✅'
	case StatusError:
		return '❌'
	default:
		return '❔'
	}
}
