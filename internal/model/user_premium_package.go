package model

import "date-apps-be/pkg/datatype"

type UserPackage struct {
	UID              string         `json:"uid"`
	UserUID          string         `json:"user_uid"`
	PremiumConfigUID string         `json:"-"`
	StartedAt        *datatype.Date `json:"started_at"`
	EndedAt          *datatype.Date `json:"ended_at"`
	Quota            int64          `json:"quota"`

	PremiumConfig *PremiumConfig `json:"premium_config"`
}

func (u *UserPackage) IsExpiredPackage() bool {
	now := datatype.NewDateNow()
	return u.EndedAt.IsBefore(now)
}
