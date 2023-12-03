package resource

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusActive  Status = "active"
	StatusError   Status = "error"
	StatusDeleted Status = "deleted"
	StatusReady   Status = "ready"
	StatusPending Status = "pending"
	StatusRunning Status = "running"

	// Job Statuses.

	StatusQueued   Status = "queued"
	StatusSucceded Status = "succeeded"
	StatusFailed   Status = "failed"
	StatusCanceled Status = "canceled"
)

func (s *Status) Emoji(resourceType Type) rune {
	if s == nil {
		return ' '
	}

	switch *s {
	case StatusActive, StatusReady, StatusSucceded:
		return '✅'
	case StatusRunning:
		if resourceType == TypeJobRun {
			return '🏃'
		}
		return '✅'
	case StatusPending, StatusQueued:
		return '🕒'
	case StatusError, StatusFailed:
		return '❌'
	case StatusDeleted, StatusCanceled:
		return '🧹'
	default:
		return '❔'
	}
}
