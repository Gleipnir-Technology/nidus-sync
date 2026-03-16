package types

import (
	"github.com/nyaruka/phonenumbers"
)

type E164 struct {
	number *phonenumbers.PhoneNumber
}

func NewE164(n *phonenumbers.PhoneNumber) *E164 {
	return &E164{
		number: n,
	}
}
func (e E164) PhoneString() string {
	return phonenumbers.Format(e.number, phonenumbers.E164)
}
