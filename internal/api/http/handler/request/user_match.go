package request

type CreateMatch struct {
	MatchUID  string `json:"match_uid" validate:"required"`
	MatchType string `json:"match_type" validate:"required"`
}
