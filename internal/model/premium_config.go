package model

type PremiumConfig struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Quota       int64  `json:"quota"`
	ExpiredDay  int64  `json:"expired_day"`
	IsActive    bool   `json:"is_active"`
}
