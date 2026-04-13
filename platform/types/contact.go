package types

import (
	"encoding/json"
	//"github.com/rs/zerolog/log"
)

type Contact struct {
	CanSMS   *bool   `db:"can_sms" json:"can_sms"`
	Email    *string `db:"email" json:"email"`
	HasEmail bool    `json:"has_email"`
	HasPhone bool    `json:"has_phone"`
	Name     *string `db:"name" json:"name"`
	Phone    *string `db:"phone" json:"phone"`
}

func (c Contact) MarshalJSON() ([]byte, error) {
	to_marshal := make(map[string]interface{}, 0)
	if c.CanSMS != nil {
		to_marshal["can_sms"] = *c.CanSMS
	}
	to_marshal["name"] = c.Name
	to_marshal["has_email"] = (c.Email != nil && *c.Email != "")
	to_marshal["has_phone"] = (c.Phone != nil && *c.Phone != "")
	//log.Debug().Msg("marshaling contact")
	return json.Marshal(to_marshal)
}
