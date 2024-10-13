package types

type ChatType string

const (
	ChatTypeDefault    ChatType = "DEFAULT"
	ChatTypeReflection ChatType = "REFLECTION"
)

// IsValid checks if the ChatType is valid
func (ct ChatType) IsValid() bool {
	switch ct {
	case ChatTypeDefault, ChatTypeReflection:
		return true
	}
	return false
}
