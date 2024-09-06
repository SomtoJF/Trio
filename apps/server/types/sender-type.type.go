package types

type SenderType string

const (
	SenderTypeUser  SenderType = "User"
	SenderTypeAgent SenderType = "Agent"
)

// IsValid checks if the SenderType is valid
func (st SenderType) IsValid() bool {
	switch st {
	case SenderTypeUser, SenderTypeAgent:
		return true
	}
	return false
}
