package risk

// Level is an ordered severity value.
type Level int

const (
	None   Level = 0
	Low    Level = 1
	Medium Level = 2
	High   Level = 3
)

// Parse converts a risk string into a Level.
func Parse(s string) Level {
	switch s {
	case "low":
		return Low
	case "medium":
		return Medium
	case "high":
		return High
	default:
		return None
	}
}

// String converts a Level back to its canonical string.
func String(l Level) string {
	switch l {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "none"
	}
}

// ShouldFail returns true when any risk in riskLevels meets or exceeds threshold.
func ShouldFail(threshold string, riskLevels []string) bool {
	if threshold == "none" {
		return false
	}
	t := Parse(threshold)
	for _, r := range riskLevels {
		if Parse(r) >= t {
			return true
		}
	}
	return false
}
