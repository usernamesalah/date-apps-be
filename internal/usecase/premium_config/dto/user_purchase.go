package dto

type UserPurchase struct {
	UserUID          string `json:"user_uid"`
	PremiumConfigUID string `json:"premium_config_uid"`
}
