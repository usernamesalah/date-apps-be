package request

type UserPurchase struct {
	PremiumConfigUID string `json:"premium_config_uid" validate:"required"`
}
