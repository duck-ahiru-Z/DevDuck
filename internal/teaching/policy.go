package teaching

import "fmt"

type Level string

const (
	Beginner     Level = "beginner"
	Intermediate Level = "intermediate"
	Advanced     Level = "advanced"
)

type Policy struct {
	Level       Level
	ShowSummary bool
	MaxHints    int
}

func DefaultPolicy() Policy {
	return PolicyForLevel(Intermediate)
}

func PolicyForLevel(level Level) Policy {
	switch level {
	case Beginner:
		return Policy{
			Level:       Beginner,
			ShowSummary: true,
			MaxHints:    3,
		}

	case Advanced:
		return Policy{
			Level:       Advanced,
			ShowSummary: false,
			MaxHints:    1,
		}

	default:
		return Policy{
			Level:       Intermediate,
			ShowSummary: true,
			MaxHints:    2,
		}
	}
}

func ParseLevel(value string) (Level, error) {
	switch Level(value) {
	case Beginner:
		return Beginner, nil

	case Intermediate:
		return Intermediate, nil

	case Advanced:
		return Advanced, nil

	default:
		return "", fmt.Errorf("unknown level: %s", value)
	}
}
