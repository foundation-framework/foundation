package session

type AccessLevel uint

const (
	AccessLevelBasic AccessLevel = iota
	AccessLevelExtended
)

var accessLevelFunDefault = func(data *Client) AccessLevel {
	if data == nil {
		return AccessLevelBasic
	}

	return AccessLevelExtended
}
