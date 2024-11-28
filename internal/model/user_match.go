package model

import (
	"date-apps-be/internal/constant"
	"date-apps-be/pkg/datatype"
)

type UserMatch struct {
	UserUID   string
	MatchUID  string
	MatchType constant.UserMatchType
	CreatedAt datatype.Time

	User  User
	Match User
}
