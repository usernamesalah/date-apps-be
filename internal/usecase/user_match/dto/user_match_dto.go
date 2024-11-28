package dto

import "date-apps-be/internal/constant"

type GetUserMatches struct {
	UserUID   string                 `json:"user_uid"`
	Page      uint64                 `json:"page"`
	Limit     uint64                 `json:"limit"`
	MatchType constant.UserMatchType `json:"match_type"`
}
